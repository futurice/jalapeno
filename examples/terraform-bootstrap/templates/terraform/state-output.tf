output "resource_group_name" {
  value = data.azurerm_resource_group.main.name
}

output "tfstate_storage_account_name" {
  value = azurerm_storage_account.tfstate.name
}

output "tfstate_storage_container_name" {
  value = azurerm_storage_container.tfstate.name
}