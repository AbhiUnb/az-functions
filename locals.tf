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
     // storage_account_access_key = ""

    }
  }





}