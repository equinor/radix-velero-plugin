# Radix Deployment Velero Plugin

This plugin is intended to assist velero in backup/restores, so that we are able to recover the state as it were in the original cluster. Information from this plugin ends up as annotations on the restored object. Currently the supported annotations are:

- equinor.com/velero-restored-status

This annotation will be picked up by the radix-operator, after the restore is done

## Building the plugins

To build the plugins, run

```bash
$ make
```

To build the image, run

```bash
$ make container
```

This builds an image tagged as `radixdev.azurecr.io/radix-velero-plugin:latest`. If you want to specify a
different name, run

```bash
$ make container IMAGE=your-repo/your-name:here
```

## Deploying the plugins

To deploy your plugin image to an Velero server:

1. Make sure your image is pushed to a registry that is accessible to your cluster's nodes.
2. Run `velero plugin add <image>`, e.g. `velero plugin add radixdev.azurecr.io/radix-velero-plugin:latest`