using namespace System.Net

param($Request, $TriggerMetadata)

# --- Auth with UAMI ---
$uamiClientId = "YOUR-UAMI-CLIENT-ID"

try {
    Connect-AzAccount -Identity -AccountId $uamiClientId | Out-Null
    Write-Output "Logged in with UAMI"
} catch {
    return @{
        statusCode = [HttpStatusCode]::Unauthorized
        body = "Failed to login using UAMI: $_"
    }
}

# --- DB Connection ---
$connectionString = $env:SQL_CONNECTION_STRING
$query = "SELECT mg_id FROM Management_Groups WHERE env_type = 'lower';"
$mgIds = @()

try {
    Add-Type -AssemblyName "System.Data"
    $connection = New-Object System.Data.SqlClient.SqlConnection $connectionString
    $connection.Open()
    $command = $connection.CreateCommand()
    $command.CommandText = $query
    $reader = $command.ExecuteReader()
    while ($reader.Read()) {
        $mgIds += $reader["mg_id"]
    }
    $reader.Close()
    $connection.Close()
} catch {
    return @{
        statusCode = [HttpStatusCode]::InternalServerError
        body = "Failed to fetch management groups: $_"
    }
}

# --- Fetch Subscriptions ---
$subs = @()
foreach ($mgId in $mgIds) {
    try {
        $subs += Get-AzManagementGroupSubscription -GroupName $mgId
    } catch {
        Write-Warning "Failed to get subs for $mgId: $_"
    }
}

# --- De-duplicate subscriptions ---
$uniqueSubs = @{}
$filteredSubs = @()
foreach ($sub in $subs) {
    if ($sub.Id -match "/subscriptions/([0-9a-fA-F-]+)$") {
        $subId = $matches[1]
        if (-not $uniqueSubs.ContainsKey($subId)) {
            $uniqueSubs[$subId] = $true
            $filteredSubs += $sub
        }
    }
}

# --- Analyze VMs ---
$finalOutput = @()

foreach ($sub in $filteredSubs) {
    if ($sub.Id -match "/subscriptions/([0-9a-fA-F-]+)$") {
        $subId = $matches[1]
        Set-AzContext -SubscriptionId $subId | Out-Null

        $vms = Get-AzVM
        foreach ($vm in $vms) {
            $vmId = $vm.Id
            $vmName = $vm.Name
            $rg = $vm.ResourceGroupName

            $metrics = Get-AzMetric -ResourceId $vmId `
                -TimeGrain ([TimeSpan]::FromMinutes(15)) `
                -StartTime (Get-Date).AddDays(-3) `
                -EndTime (Get-Date) `
                -MetricName "Percentage CPU" `
                -Aggregation Average

            # Group by day
            $cpuByDay = $metrics.Data | Group-Object { $_.TimeStamp.Date }
            foreach ($group in $cpuByDay) {
                $cpuValues = $group.Group | Where-Object { $_.Average -ne $null } | Select-Object -ExpandProperty Average
                if (-not $cpuValues) { continue }

                $avg = ($cpuValues | Measure-Object -Average).Average
                $std = ($cpuValues | Measure-Object -StandardDeviation).StandardDeviation
                $threshold = 10

                # Find idle slots
                $timestamps = $group.Group | Sort-Object TimeStamp
                $idleStreak = @()
                $longestIdle = @()
                $spikeFound = $false

                foreach ($entry in $timestamps) {
                    if ($entry.Average -lt $threshold) {
                        $idleStreak += $entry
                        if ($idleStreak.Count -ge 12) {
                            $longestIdle = $idleStreak
                        }
                    } else {
                        $spikeFound = $true
                        $idleStreak = @()
                    }
                }

                $stopTime = $startTime = $null
                if ($longestIdle) {
                    $stopTime = $longestIdle[0].TimeStamp.ToString("HH:mm")
                    $nextUsage = $timestamps | Where-Object { $_.TimeStamp -gt $longestIdle[-1].TimeStamp -and $_.Average -ge $threshold }
                    if ($nextUsage) {
                        $startTime = ($nextUsage[0].TimeStamp).AddMinutes(-60).ToString("HH:mm")
                    }
                }

                $finalOutput += [PSCustomObject]@{
                    Subscription = $subId
                    ResourceGroup = $rg
                    VMName = $vmName
                    Date = $group.Name.ToString("yyyy-MM-dd")
                    StopTime = $stopTime
                    StartTime = $startTime
                    Status = if ($stopTime) { "Idle detected" } else { "No idle window" }
                }
            }
        }
    }
}

# --- Return JSON response ---
return @{
    statusCode = [HttpStatusCode]::OK
    body = ($finalOutput | ConvertTo-Json -Depth 4)
    headers = @{ "Content-Type" = "application/json" }
}
