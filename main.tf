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

