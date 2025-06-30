variable "cc_location" {
  type        = string
  description = "Canada Central Region"
}

variable "cc_core_resource_group_name" {
  type        = string
  description = "Resource Group Name for McCain Foods Manufacturing Digital Shared Azure Components in Canada Central"
}


variable "MF_DM_CC_CORE_appSP_Name" {
  description = "Web App Service Plan name for McCain Foods MF Digital Canada Central"
  type        = string
}

variable "MF_DM_CC_CORE_Webapp_Name" {
  description = "Web App name for McCain Foods MF Digital Canada Central"
  type        = string
}

variable "cc_core_function_apps" {
  type        = map(any)
  description = "Map of Function Apps with their configuration"
}