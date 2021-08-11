# Radix Velero Plugin

This plugin is intended to assist velero in backup/restores, so that we are able to recover the state as it were in the original cluster. Information from this plugin ends up as annotations on the restored object. Currently the supported annotations are:

- equinor.com/velero-restored-status

This annotation will be picked up by the radix-operator, after the restore is done

## Building the plugins

Docker images are automatically build and pushed to both radixdev and radixprod ACR. 
A push to the master branch builds a new image with tag **master-latest**, and push to release branch tags the image with **release-latest**

## Plugin deployment

To deploy your plugin image to an Velero server, there are two options.

### Manual deployment

1. Make sure your image is pushed to a registry that is accessible to your cluster's nodes.
2. Run `velero plugin add <image>`, e.g. `velero plugin add radixdev.azurecr.io/radix-velero-plugin:master-latest`

### Deployment using Flux

The regular way we deploy to cluster is using our [Flux repo](https://github.com/equinor/radix-flux).
What the velero CLI does is to add the plugin as an init container to the deployment. We do the same thing in the the Flux repo, by
adding it to the Helm release.

Development clusters should use the **master-latest** tag, and plauground and production clusters should use the **release-latest** tag.
