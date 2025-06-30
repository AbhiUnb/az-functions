output "cc_location" {
  value       = var.cc_location
  description = "Canada Central Region"
}

output "cc_core_resource_group_name" {
  value       = var.cc_core_resource_group_name
  description = "Resource Group Name for McCain Foods Manufacturing Digital Shared Azure Components in Canada Central"
}

output "MF_DM_CC_CORE_appSP_Name" {
  value       = var.MF_DM_CC_CORE_appSP_Name
  description = "Web App Service Plan name for McCain Foods MF Digital Canada Central"
}

output "MF_DM_CC_CORE_Webapp_Name" {
  value       = var.MF_DM_CC_CORE_Webapp_Name
  description = "Web App name for McCain Foods MF Digital Canada Central"
}

# output "function_app_url" {
#   value       = module.avm-res-web-site["MF-MDI-CC-GHPROD-DDDS-AFUNC"].default_hostname
#   description = "The default hostname of the Azure Function App"
# }





output "function_app_url" {
  value       = azurerm_linux_function_app.my_function.default_hostname
  description = "The default hostname of the Azure Function App"
}

output "function_app_sku" {
  value       = azurerm_service_plan.MFDMCCASPAFUNC.sku_name
  description = "The SKU name of the Function App Service Plan"
}

output "function_app_runtime" {
  value       = azurerm_linux_function_app.my_function.site_config[0].application_stack[0].node_version
  description = "The runtime version of the Function App"
}

# output "function_app_timeout" {
#   value       = azurerm_linux_function_app.my_function.site_config[0].function_timeout
#   description = "The function app timeout setting"
# }
output "function_app_deployment_source" {
  value = "GitHub Actions" # or "ZipDeploy" if you use func publish
}