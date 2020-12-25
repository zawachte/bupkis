# bupkis

container registries are bupkis ... view images and stats from the commandline.


## install cli

install from the latest [release artifacts](https://github.com/zwachtel11/peg/releases):

* Linux

  ```sh
  curl -LO https://github.com/zawachte-msft/bupkis/releases/download/v0.0.1/bupkis
  mv bupkis /usr/local/bin/
  ```

* macOS

  ```sh
  curl -LO https://github.com/zwachtel11/peg/releases/download/v0.0.1/bupkis-darwin
  mv bupkis-darwin /usr/local/bin/bupkis
  ```

* Windows

  Add `%USERPROFILE%\bin\` to your `PATH` environment variable so that `peg.exe` can be found.
  ```sh
  curl.exe -sLO  https://github.com/zwachtel11/bupkis/releases/download/v0.0.1/bupkis.exe
  copy bupkis.exe %USERPROFILE%\bin\
  set PATH=%USERPROFILE%\bin\;%PATH%

## usage

First use bupkis to login into your private OCI compliant container registry.

```sh
bupkis login bupkisimages.azurecr.io -u <CR_USERNAME> -p <CR_PASSWORD>
```

To search all of the private registries which you have access, login to all of them with either `bupkis` or `docker` cli. 


Then to list all the images in a private container registry run the following.

```
bupkis list bupkisimages.azurecr.io
```

If you want to list the images in all of the private container registries you can run.

```
bupkis list
```

If you just want to see all of the tags for a single image you can run.

```
bupkis get bupkisimages.azurecr.io/docs-image
```

Or with a tag to just make sure it is pushed up successifully.

```
bupkis get bupkisimages.azurecr.io/docs-image:latest
```

## roadmap
