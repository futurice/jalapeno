{{ $environments := splitList "," .Variables.ENVIRONMENTS }}
{{ $resource_groups := splitList "," .Variables.RESOURCE_GROUP_NAMES }}

locals {
  resource_groups = {
    {{- range $index, $env := $environments }}
    {{ $env | quote }}: {{ (index $resource_groups $index) | quote }}
    {{- end }}
    "default": {{ (index $resource_groups 0) | quote }}
  }
}

data "azurerm_client_config" "current" {
}

data "azurerm_resource_group" "main" {
  name     = local.resource_groups[terraform.workspace]
}

resource "azurerm_storage_account" "tfstate" {
  name                     = "{{ (printf "tfs%.11s%.6s" (regexReplaceAll "[^a-z0-9]" (.Variables.SERVICE_NAME | lower) "") (sha1sum .Recipe.Anchor)) }}${terraform.workspace}"
  resource_group_name      = data.azurerm_resource_group.main.name
  location                 = data.azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_2"
}

#
# TODO: family of stableRandomX helper functions for sprig where they always give the same value for the same id, e.g. "gimme 6 random alphanumeric characters that don't change on recipe upgrade"

resource "azurerm_storage_container" "tfstate" {
  name                  = "tfstate"
  storage_account_name  = azurerm_storage_account.tfstate.name
}

resource "azurerm_role_assignment" "tfstate" {
  scope                = azurerm_storage_container.tfstate.resource_manager_id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "local_file" "backend_config" {
  filename = "backend.tf"
  content = <<-EOT
terraform {
	backend "azurerm" {
    use_azuread_auth = true
	}
}
EOT
}
