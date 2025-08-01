param($Request)
# ---AUthentication using UAMI-----
$uamiClientId = "c1d0570e-7874-4f47-a752-14222f36dc00"

try {
    Connect-AzAccount -Identity -AccountId $uamiClientId | Out-Null
    Write-Output "Logged in with User Assigned Managed Identity (ClientId: ${uamiClientId})"
} catch {
    Write-Error "Failed to login with User Assigned Managed Identity: $_"
    return
}

# --- FETCH MANAGEMENT GROUP IDS FROM AZURE SQL DATABASE ---
$connectionString = ${env:SQL_CONNECTION_STRING}
$query = 'SELECT mg_id FROM Management_Groups WHERE env_type = ''lower'';'

try {
    Add-Type -AssemblyName "System.Data"

    $connection = New-Object System.Data.SqlClient.SqlConnection
    $connection.ConnectionString = $connectionString
    $connection.Open()

    $command = $connection.CreateCommand()
    $command.CommandText = $query

    $reader = $command.ExecuteReader()
    $managementGroupIds = @()
    while ($reader.Read()) {
        $managementGroupIds += $reader["mg_id"]
    }

    $reader.Close()
    $connection.Close()

    Write-Output "Fetched Management Group IDs from DB: $($managementGroupIds -join ', ')"
} catch {
    Write-Error "Failed to fetch Management Group IDs from DB: $_"
    return
}

# --- GET ALL SUBSCRIPTIONS UNDER NON-PROD MGs ---
$allSubscriptions = @()

foreach ($mgId in $managementGroupIds) {
    try {
        Write-Output "Getting subscriptions under Management Group: ${mgId}"
        $subs = Get-AzManagementGroupSubscription -GroupName $mgId
        if ($subs) {
            $allSubscriptions += $subs
            Write-Output "Found $($subs.Count) subscriptions under MG ${mgId}"
        } else {
            Write-Output "No subscriptions found under MG ${mgId}"
        }
    } catch {
        Write-Error "Failed to get subscriptions for MG ${mgId}: $_"
    }
}

# Remove duplicates based on Subscription ID
$uniqueSubs = @{}
$filteredSubscriptions = @()

foreach ($sub in $allSubscriptions) {
    if ($sub.Id -match "/subscriptions/([0-9a-fA-F-]+)$") {
        $subId = $matches[1]

        if (-not $uniqueSubs.ContainsKey($subId)) {
            $uniqueSubs[$subId] = $true
            $filteredSubscriptions += $sub
        }
    } else {
        Write-Warning "Cannot extract subscription ID from $($sub.Id)"
    }
}

$allSubscriptions = $filteredSubscriptions
Write-Output "Total unique subscriptions to process: $($allSubscriptions.Count)"

# --- FOR EACH SUBSCRIPTION, FIND VMs AND APPLY PARKING LOGIC ---
foreach ($sub in $allSubscriptions) {
    $fullSubId = $sub.Id
    $subName = $sub.Name

    if ($fullSubId -match "/subscriptions/([0-9a-fA-F-]+)$") {
        $subId = $matches[1]
    } else {
        Write-Error "Cannot extract subscription ID from $fullSubId"
        continue
    }
    
    Write-Output "Processing subscription: ${subName} (${subId})"

    try {
        $subDetails = Get-AzSubscription -SubscriptionId $subId
        $tenantId = $subDetails.TenantId
    } catch {
        Write-Error "Failed to get subscription details for $subId : $_"
        continue
    }

    

    try {
        Set-AzContext -SubscriptionId $subId -TenantId $tenantId | Out-Null
        Write-Output "Context set for Subscription: $subName with TenantId: $tenantId"
    } catch {
        Write-Error "Failed to set context for ${subName} (${subId}) with TenantId: $tenantId : $_"
        continue
    }

    try {
       $vms = Get-AzVM
       foreach ($vm in $vms) {
         $vmId = $vm.Id
         $vmName = $vm.Name
         Write-Output "Analyzing VM- $vmName"

         $end = Get-Date
         $start = $end.AddDays(-7)

         $metrics = Get-AzMetric -ResourceId $vmId -TimeGrain 00:15:00 -StartTime $start -EndTime $end -MetricName "Percentage CPU", "OS Disk Read Bytes/Sec", "OS Disk Write Bytes/Sec"

         # Extract CPU data points and timestamps
         $cpuPoints = $metrics | Where-Object { $_.MetricName.Value -eq "Percentage CPU" } | ForEach-Object { $_.Data }
         $cpuAverages = $cpuPoints | Where-Object { $_.Average -ne $null } | Select-Object -ExpandProperty Average
         $cpuTimestamps = $cpuPoints | Where-Object { $_.Average -ne $null } | Select-Object -ExpandProperty TimeStamp

         $windowSize = 12  # 3 hours of 15-min data
         $foundIdle = $false
        
         # Find 3-hour idle window (CPU <= 10%)
         for ($i = 0; $i -le ($cpuAverages.Count - $windowSize); $i++) {
             $window = $cpuAverages[$i..($i + $windowSize - 1)]
             if ($window -and ($window | Where-Object { $_ -gt 10 }) -eq $null) {
                 $stopTime = $cpuTimestamps[$i]
                 Write-Output "  ↳ Stop time suggestion for $vmName  $stopTime"
                 $foundIdle = $true
                 break
             }
         }

        if (-not $foundIdle) {
            Write-Output "  ↳ No 3-hour idle window found for $vmName"
        }

        # # Find last time CPU > 20% (suggested start time)
        # for ($j = $cpuAverages.Count - 1; $j -ge 1; $j--) {
        #     if ($cpuAverages[$j] -gt 20) {
        #         $index = [Math]::Max(0, $j - 4)  # 1 hour = 4 intervals before
        #         $startTime = $cpuTimestamps[$index]
        #         Write-Output "  ↳ Start time suggestion for $vmName: $startTime"
        #         break
        #     }
        # }
    }
} catch {
    Write-Error "Error processing VMs: $_"
}
}

Push-OutBinding -Name Response -Value @{
    StatusCode = 200
    Body = ($vmName | ConvertTo-Json -Depth 4)
}
return $Response
