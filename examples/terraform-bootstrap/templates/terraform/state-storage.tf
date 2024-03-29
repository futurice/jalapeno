{{- define "rg_block_type" -}}
{{- ternary "resource" "data" .Variables.CREATE_RESOURCE_GROUPS -}}
{{- end -}}

{{- define "resource_tag" -}}
{{- printf "%.6s" (sha1sum .ID) -}}
{{- end -}}

{{- define "storage_account_name_prefix" -}}
{{- printf "tfs%.11s" (regexReplaceAll "[^a-z0-9]" (.Variables.SERVICE_NAME | lower) "") -}}
{{- end -}}

locals {
  resource_groups = {
    {{- range $index, $env := .Variables.ENVIRONMENTS }}
    {{ $env.NAME | quote }} : {{ $env.RESOURCE_GROUP_NAME | quote }}
    {{- end }}
    "default" : {{ (index .Variables.ENVIRONMENTS 0).RESOURCE_GROUP_NAME | quote }}
  }

  resource_tag = "{{ template "resource_tag" . }}"
}

data "azurerm_client_config" "current" {
}

{{ if .Variables.CREATE_RESOURCE_GROUPS -}}
resource "azurerm_resource_group" "main" {
  name     = local.resource_groups[terraform.workspace]
  location = {{ .Variables.RESOURCE_GROUP_LOCATION | quote }}
}
{{- else -}}
data "azurerm_resource_group" "main" {
  name = local.resource_groups[terraform.workspace]
}
{{- end }}

resource "azurerm_storage_account" "tfstate" {
  name                     = "{{ template "storage_account_name_prefix" . }}{{ template "resource_tag" . }}${terraform.workspace}"
  resource_group_name      = {{ template "rg_block_type" . }}.azurerm_resource_group.main.name
  location                 = {{ template "rg_block_type" . }}.azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_2"
}

resource "azurerm_storage_container" "tfstate" {
  name                 = "tfstate"
  storage_account_name = azurerm_storage_account.tfstate.name
}

resource "azurerm_role_assignment" "tfstate" {
  scope                = azurerm_storage_container.tfstate.resource_manager_id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "local_file" "backend_config" {
  filename = "backend.tf"
  content  = <<-EOT
terraform {
  backend "azurerm" {
    use_azuread_auth = true
  }
}
EOT
}
