# wplite

serverless WordPress. need I say more?

## Dependencies

Before you can use `wplite`, you must have the following dependencies installed:

- `docker`
- `git`
- `git-lfs`

## How does this work?

`wplite` intentionally does not try to "reinvent the wheel" in the WordPress development space. By relying on core WordPress constructs, most WordPress plugins and themes will work out of the box with `wplite`. This is a key feature of `wplite` - it is not a new CMS, it is a new way to run WordPress.

`wplite` currently relies on a fork of upstream WordPress core which adds support for sqlite. This enables the CMS to run in a Docker container without the need for a MySQL server, and the entire database can be committed to git and shared between developers. Once this fork is merged into upstream WordPress, `wplite` will be able to use the official WordPress core.

`wplite` also includes a static site generator to generate files which can then be hosted on any static site hosting service.

From there, `wplite` simply relies on `git` and `git-lfs` to store the entire state of the site in whichever `git`-compatible VCS and workflow you prefer.

## Usage

### Creating a Site

Create an empty directory for your new site, and navigate to it:

```sh
mkdir my-site
cd my-site
```

`wplite` uses a `.wplite-env` file in the current directory to configure the WordPress site. This file is a simple key-value pair file, with each line in the format `KEY=VALUE`. The following keys are supported:

```sh
WP_TITLE=my-site
WP_USER=admin
WP_PASS=password
WP_EMAIL=hello@example.com
WP_PORT=80
WP_THEME=twentytwentyfour
```

You can create this file manually, or let `wplite` prompt you for the values when you first run `wplite start` in an empty directory.

On first start, `wplite` will create a `wp-content` directory and a `.htaccess` file in the current directory. The first start of a site can take 1-2 minutes to get everything set up. Once the site is initialized, your browser will open to the CMS. You can then log in with the username and password you specified in `.wplite-env`.

### Starting an Existing Site

If you have an existing site directory, you can simply run `wplite start` in that directory to start the CMS. Once the site is initialized (usually takes about 30-45 seconds), your browser will open to the CMS. You can then log in with the username and password you specified in `.wplite-env`.

Note that since a `wplite` site is simply a deployment of WordPress running in a Docker container, `wplite start` is effectively a wrapper around:

```bash
docker run -d -p 80:80 \
    -v $(pwd)/wp-content:/var/www/html/wp-content \
    -v $(pwd)/.htaccess:/var/www/html/.htaccess \
    --env-file .wplite-env \
    --name wplite \
    robertlestak/wplite:latest
```

Knowing this will come in handy when we discuss deployment models below.

### Building the Site

Once you have made changes to the site, you can run `wplite build` in the directory to generate the static site contents in `wp-content/static`. This will generate a static version of the site, which can be hosted on any static site hosting service. You can validate this with something like [serve](https://www.npmjs.com/package/serve):

```sh
npx serve wp-content/static
```

By default, running `wplite build` will stop the local CMS on completion. You can use the `-no-stop` flag to keep the CMS running after the build.

If you do not have the CMS running locally when you run `wplite build`, it will start the CMS, build the static site contents, and then stop the CMS.

### Stopping the Site

You can stop the CMS by running `wplite stop` in the site directory. This will stop the CMS and free up the port and resources for other applications.

**NOTE**: You should stop the site before committing changes to git, as the `.ht.sqlite` file will be locked while the CMS is running.

## Deployment Models

`wplite` was designed to provide an alternative to traditional WordPress hosting. WordPress itself is not very scalable, but the plugins, themes, and WYSIWYG editor make it a popular choice for the majority of sites on the internet, both small and extremely large.

For small sites, the overhead of an always-on server, a database, and a CMS is an unnecessary expense that historically has been difficult to avoid when using WordPress.

At scale, WordPress exhibits a number of problems which generally necessitate a more complex hosting setup including horizontal and vertical autoscaling, caching, and a CDN. This adds many layers of complexity to the hosting setup, and generally requires a team of DevOps engineers to manage.

`wplite` aims to provide a more slimmed-down implementation of WordPress which can be entirely hosted on a static site hosting service. This provides a number of benefits:

- **Cost**: static site hosting is generally much cheaper than traditional hosting, and can often be free for small sites.
- **Scalability**: static site hosting services are generally designed to handle large amounts of traffic, and can be scaled up and down as necessary.
- **Security**: static site hosting services are generally more secure than traditional hosting, and are designed to handle DDoS attacks and other security threats.
- **Simplicity**: static site hosting services are generally easier to set up and manage than traditional hosting, and require less maintenance.
- **Portability**: static site hosting services are generally more portable than traditional hosting, and can be easily moved between providers.
- **Performance**: static site hosting services are generally faster than traditional hosting, and can be easily cached and distributed globally.
- **Reliability**: static site hosting services are generally more reliable than traditional hosting, and have better uptime guarantees.
- **Flexibility**: static site hosting services are generally more flexible than traditional hosting, and can be easily integrated with other services.

This does come with a trade-off, however. Since the CMS is not running on the server, you cannot use server-side code in your site. This means you cannot use PHP, and you cannot use server-side databases. This also means you cannot use server-side forms, comments, or other dynamic content. However, you can use third-party services to handle these features, and you can use a plugin to handle comments and form submissions, and call the service over AJAX.

Generally these purpose-built microservice / Saas solutions are more scalable and secure than a traditional WordPress setup, and are a better choice for most sites both small and large.

### Static Site Workflow

`wplite` makes WordPress site development and CMS management just like any other `git` project. The typical workflow is:

- pull latest changes from upstream (`git pull`)
- start CMS locally (`wplite start`)
- make changes to site (manual in CMS)
- stop CMS (`wplite stop`)
- push changes to remote (`git push`)
    - let CI/CD build static site contents (`wplite build`)
    - let CI/CD publish static site contents to hosting service

### Containerized Workflow

If your site requires server-side code, you can deploy `wplite` as a single container on a server. This is a more traditional WordPress deployment model, but it still has some benefits over a traditional WordPress deployment, namely the reduced overhead of a database server, and the ability to scale to zero when not in use. However since the `sqlite` database will be accessed from within the container, you will not be able to horizontally scale the CMS, and will need to rely on vertical scaling to handle increased traffic.

In this model it is still strongly recommended to offload as much static content as possible to a static service / CDN, and only pass through to the CMS for dynamic content.

## FAQ

### WordPress version selection

Currently, `wplite` only works on a fork on the latest version of WordPress (`6.5-alpha`). Once this fork is merged into upstream WordPress, `wplite` will be able to use the official WordPress core and have the same version selection as upstream WordPress. If this merge takes longer than expected, newer versions of WordPress will be supported by `wplite` as they are released.

Once additional version support is added, you will be able to specify the version of WordPress to use in the `.wplite-env` file.

### Does `wplite` support plugin X or theme Y?

It depends. If they followed semantic WordPress development practices, then yes, they should work out of the box. If they are doing something non-standard, then they may not work. If you have a specific plugin or theme you would like to use, you can try it out and see if it works. If it doesn't, simply roll back your local changes - the beauty of `git` based WordPress development!

### Will comments and contact forms work?

By default, no. `wplite build` generates a completely static site, and comments and contact forms require a server to handle the form submissions. However, you can use a third-party service to handle form submissions, and you can use a third-party service to handle comments. You can also use a plugin to handle comments and form submissions, and call the service over AJAX.

### Why is there no `publish` command?

You may notice the `wplite` CLI does _not_ have a `publish` command. This was an intentional omission. `wplite` is responsible for managing the CMS and the static site generator. Once it generates the static site contents in `wp-content/static`, you can take those files and publish them to ANY static service - AWS S3, Netlify, Vercel, a local FTP server only accessible over some arcane VPN, etc.

If `wplite` included a `publish` command, it would start to become a deployment tool, and that's not the goal of this project. The goal is to provide a lightweight, flexible, and portable CMS and static site generator, and let you decide how to publish the static site contents.

While this may seem like a cop-out, it's actually a feature. It means you can use `wplite` with any static site hosting service, and you're not locked into a specific deployment platform or workflow.

If you're not happy with this answer, check out the `examples/` directory for some example end-to-end CI/CD workflows.

### Committing media files to git?

By default, `wplite` will assume you have `git-lfs` configured to handle media files. This is the easiest way to handle media files for smaller projects. However you can also use a plugin to offload media files to a cloud storage service directly, and in those cases, you would not need to use `git-lfs`. Since that is a more advanced WordPress-specific setup, it is not covered in this README.

### Committing `.wplite-env` and `.ht.sqlite` to git?

It is generally considered best practice to not commit `.env` files to git, as they often contain sensitive information. However, in the case of `.wplite-env`, it is necessary to commit this file to git, as it contains the environment variables necessary to run the CMS and static site generator. While you are technically specifying a username and password to access the CMS, since all development is done locally, and the CMS is not accessible from the internet, this is not a security risk, and generally these are set to simple values like `admin` and `password`.

`.ht.sqlite` is a SQLite database file that contains the CMS data. It is also necessary to commit this file to git, as it contains the data necessary to run the CMS and static site generator and share between developers. Similar to above, while general best practice dictates not to commit database files to git, in this case it is necessary to do so, as all of the site data is stored in this file.

If your site is expected to handle user inputs or sensitive information, you can still use `wplite` to build the user interface and static site, but you should use something more scalable and secure for the database and call the service over AJAX.

## Limitations

- `wplite` only works on a fork on the latest version of WordPress (`6.5-alpha`). Once this fork is merged into upstream WordPress, `wplite` will be able to use the official WordPress core and have the same version selection as upstream WordPress. If this merge takes longer than expected, newer versions of WordPress will be supported by `wplite` as they are released.
- `wplite` only works on Linux and macOS. Windows support is not planned at this time.
  - _technically_ it could work on Windows all the same, but since Docker support on Windows requires WSL and most Windows users don't have WSL set up, it's not worth the effort to support Windows at this time.
- Currently, only port `80` is supported. If you set the `WP_PORT` to anything other than `80`, while the CMS will load, it will not be able to generate the static site contents. This is a known limitation and will be fixed in a future version of `wplite`.
- For some a `git`-based workflow may be a limitation. If you are not familiar with `git`, you will need to learn the basic concepts (push, pull, etc) to use `wplite`.
- `wplite` is intended for small sites with a relatively low number of editors, or at least editors whose work does not overlap. Since the entire state of the database is stored in the `.ht.sqlite` file, if two editors make changes to the site at the same time, the changes will conflict and one will be lost. `git` doesn't do a great job of handling binary files, so it's not possible to merge the changes. A future version of `wplite` may further decouple changes to the database, but for now, it is recommended to only have one editor working on the site at a time.