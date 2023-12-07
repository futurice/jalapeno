# Terraform Bootstrap Example recipe

This recipe demonstrates how to bootstrap Terraform state management in Azure.
The goal is to have separate state storage for a number of different
environments, all on Azure, and to manage the state storage itself using
Terraform. This is a common pattern in IaC projects.

The recipe is mainly intended to be used together with GitHub Actions, but that
functionality is optional and managed by a user setting that is set when
executing the recipe. If GitHub Actions is used, pipeline files are generated
such that IaC changes flow through a format check, validate, plan and apply
pipeline.

## Prerequisites

Pre-creating resources is optional, but permissions management is easier if you
do. The following resources are required:

1. A subscription
2. As many resource groups as you want to have environments (e.g. dev, qa,
   prod). These resource groups should be empty, but they don't have to be.
3. A service principal with contributor access to the resource groups, and a
   client secret for the SP.

### Generating a service principal

1. Go to Azure Portal and the Entra ID blade and add an application. It does not
   matter what the redirect URL is. Everything can be left at default.
2. For each resource group, go to the IAM blade and add the application as a
   contributor.
3. Go to the Certificates & secrets blade for the Service Principal and add a
   client secret. Copy the secret value.

## Usage

Authenticate to Azure:

```shell
az login
```

You can use either your own account (if you have the necessary permissions on
the target subscription) of the Service Principal generated earlier.

Run the following commands:

```shell
cd terraform && task init
```

If you chose to generate a CI/CD pipeline, the init task will prompt for the
subscription ID, tenant ID, client ID and client secret for each of the
environments. Use the values for the Service Principal generated earlier.
These values will be stored as secrets in GitHub.

If you used the Service Principal credentials when running `task init`, you are
done. If not, you need to also assign the "Storage Blob Data Contributor" role
to the Service Principal on the storage accounts created by the recipe. You can
do this in the Azure Portal in IAM blades of the resource groups.
