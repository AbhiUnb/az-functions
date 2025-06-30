cc_location = "canadacentral"

cc_core_resource_group_name = "rg-mccain-core-prod"

MF_DM_CC_CORE_appSP_Name = "MFDMCCPRODASPAFUNCie"

MF_DM_CC_CORE_Webapp_Name = "MFDMCCPRODFUNCTIONAPP"

cc_core_function_apps = {
  myfunctionapp = {
    name                        = "MFDMCCPRODFUNCTIONAPP"
    location                    = "canadacentral"
    os_type                     = "Linux"
    storage_account_name        = "mfmdiccprodsa"
    storage_account_access_key  = "your-storage-account-access-key"
    storage_account_rg          = "rg-mccain-core-prod"
    network_name                = "mccain-vnet-prod"
    subnet_name                 = "function-subnet"
    user_assigned_identity_name = "mccain-func-identity"
    user_assigned_identity_rg   = "rg-mccain-core-prod"
    app_insights_name           = "mccain-func-appinsights"
    app_insights_rg             = "rg-mccain-core-prod"
    key_vault_name              = "mccain-keyvault-prod"

    additional_app_settings = {
      FUNCTIONS_WORKER_RUNTIME = "node"
      WEBSITE_RUN_FROM_PACKAGE = "https://storageaccount.blob.core.windows.net/container/package.zip?<sas_token>"
      # Example Key Vault reference syntax
      mySecretSetting = "@Microsoft.KeyVault(SecretUri=https://mccain-keyvault-prod.vault.azure.net/secrets/mySecret/)"
    }
  }
}


tfvars------

# # resource "azurerm_app_service_plan" "MFDMCCASPAFUNC" {
# #   name                = "MFDMCCPRODASPAFUNC"
# #   resource_group_name = var.cc_core_resource_group_name
# #   location            = var.cc_location
# #   kind                = "FunctionApp"
# #   sku {
# #     tier = "PremiumV2"
# #     size = "P1v2"
# #   }
# # }


resource "azurerm_service_plan" "MFDMCCASPAFUNC" {
  resource_group_name = var.cc_core_resource_group_name
  location            = var.cc_location
  name                = "MFDMCCPRODASPAFUNC"
  os_type             = "Linux"
  sku_name            = "Y1"
  tags                = local.tag_list_1
}

# module "avm-res-web-site" {
#   source              = "Azure/avm-res-web-site/azurerm"
#   version             = "0.16.4"
#   for_each            = local.functionapp
#   name                = each.value.name
#   resource_group_name = var.cc_core_resource_group_name
#   location            = var.cc_location

#   kind = each.value.kind

#   # Uses an existing app service plan
#   os_type                  = azurerm_service_plan.MFDMCCASPAFUNC.os_type
#   service_plan_resource_id = azurerm_service_plan.MFDMCCASPAFUNC.id

#   # Uses an existing storage account
#   storage_account_name       = each.value.storage_account_name
#   storage_account_access_key = each.value.storage_account_access_key
#   # storage_uses_managed_identity = true
#   site_config = {
#     always_on = false
#   }

#   tags = local.tag_list_1


# }

resource "azurerm_linux_function_app" "my_function" {
  name                = "my-node-function-app1"
  resource_group_name = var.cc_core_resource_group_name
  location            = var.cc_location
  service_plan_id     = azurerm_service_plan.MFDMCCASPAFUNC.id

  storage_account_name       = "mfmdiccprodfunctionsa"
  storage_account_access_key = "xZLtjw2G2FngjwNTihWIASyGpNz/NdzkNtWshF7jpzhFxwwO7kxNX2Gg5vRjgIv9pebBx+6N3Xi4+AStUtKW7A=="

  site_config {
    always_on = false

    application_stack {
      node_version = "18"
    }

    //function_app_timeout = "PT5M"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = local.tag_list_1
}

main.tf





locla.tf


locals {
  tag_list_1 = {
    "Application Name" = "McCain DevSecOps"
    "GL Code"          = "N/A"
    "Environment"      = "sandbox"
    "IT Owner"         = "mccain-azurecontributor@mccain.ca"
    "Onboard Date"     = "12/19/2024"
    "Modified Date"    = "N/A"
    "Organization"     = "McCain Foods Limited"
    "Business Owner"   = "trilok.tater@mccain.ca"
    "Implemented by"   = "trilok.tater@mccain.ca"
    "Resource Owner"   = "trilok.tater@mccain.ca"
    "Resource Posture" = "Private"
    "Resource Type"    = "Terraform POC"
    "Built Using"      = "Terraform"
  }


  functionapp = {
    MF-MDI-CC-GHPROD-DDDS-AFUNC = {
      name                       = "MF-MDI-CC-GHPROD-DDDS-AFUNCie"
      kind                       = "functionapp"
      storage_account_name       = "mfmdiccprodfunctionsa"
      storage_account_access_key = "xZLtjw2G2FngjwNTihWIASyGpNz/NdzkNtWshF7jpzhFxwwO7kxNX2Gg5vRjgIv9pebBx+6N3Xi4+AStUtKW7A=="

    }
  }





}



