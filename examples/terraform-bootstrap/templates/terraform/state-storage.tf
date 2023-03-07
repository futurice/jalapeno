data "azurerm_resource_group" "main" {
  name     = {{ quote .Variables.RESOURCE_GROUP_NAME }}
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
# TODO: "anchor" / "id" in rendered recipe (generated when executed, somehow uniquely identifies recipe + project combo)
# TODO: family of stableRandomX helper functions for sprig where they always give the same value for the same id, e.g. "gimme 6 random alphanumeric characters that don't change on recipe upgrade"

resource "azurerm_storage_container" "tfstate" {
  name                  = "tfstate"
  storage_account_name  = azurerm_storage_account.tfstate.name
}

resource "local_file" "backend_config" {
  filename = "backend.tf"
  content = <<-EOT
terraform {
	backend "azurerm" {
		storage_account_name = "${azurerm_storage_account.tfstate.name}"
		container_name = "${azurerm_storage_container.tfstate.name}"
	}
}
EOT
}

output "tfstate_storage_account_name" {
  value = azurerm_storage_account.tfstate.name
}