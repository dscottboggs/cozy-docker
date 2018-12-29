# cozy-docker

[![Docker Hub stars](https://img.shields.io/docker/stars/moritzheiber/cozy-stack.svg)](https://hub.docker.com/r/moritzheiber/cozy-stack) [![Docker Hub pulls](https://img.shields.io/docker/pulls/moritzheiber/cozy-stack.svg)](https://hub.docker.com/r/moritzheiber/cozy-stack) [![Docker Hub automated](https://img.shields.io/docker/automated/moritzheiber/cozy-stack.svg)](https://hub.docker.com/r/moritzheiber/cozy-stack)

A small(ish) image to run the latest version of the [Cozy cloud (v3)](https://cozy.io).

## Deployment:

1. Configure your DNS and reverse-proxy to point at your domain/subdomain that
   you're using for each instance, and each of these subdomains on that domain
   or subdomain. For simplicity's sake we're going to refer to the rest of this
   deployment as though it were being deployed under the domain *example.com*.
   With this in mind (that you'll have to replace *example.com* with the domain
   or subdomain for this instance/user), these are all the domains and
   subdomains you need to point at this stage.
 - example.com
 - onboarding.example.com
 - settings.example.com
 - drive.example.com
 - photos.example.com
 - home.example.com
 - store.example.com
 - collect.example.com
   If you use Traefik, that means you need to add the following labels to the
   container under the `cozy` heading in `docker-compose.yml`.
 - traefik.docker.network: <the network traefik is configured to watch>
 - traefik.cozy_cloud.frontend.rule: Host:<each of the above domains separated by commas>
 - traefik.cozy_cloud.port: 8080
 - traefik.cozy_cloud.protocol: http
 - traefik.enable: "true"
   In addition, if you install other apps, you'll need to add that app's
   subdomain to your reverse proxy configuration.

   For DNS, it's just a matter of adding *example.com* and *\*.example.com* to
   to your A-Records, pointed at your reverse proxy.

2. Generate a new long random password using alphanumeric characters, and store
   it somewhere safe, like in the project root and named `cozy-admin-passphrase`
   (which is in the `.gitignore` file), right next to `cozy.yml`. Encrypt it with
   ```sh
   passphrase=$(cat cozy-admin-passphrase) go run scripts/encrypt_pw.go
   ```
   from the project root (or an scrypt generator of your choice), and place it
   at (or redirect the command output to) `mounts/cozy-conf/hashed-cozy-admin-passphrase`.

3. Run `docker-compose up --build -d && docker-compose logs -f`. This will allow
   allow you to view the logs without stopping everything when you exit the logs
   view. You'll need to watch the logs because the main application will stop
   after the first run because CouchDB won't be up yet. Once CouchDB is up
   (youll start seeing its log entries), press ctrl+C to exit the logview and
   run `docker-compose up -d` again. Your Cozy stack should be running now.

4. Create an instance: run the following command to create a Cozy instance
    ```sh
    docker-compose exec -e "COZY_ADMIN_PASSWORD=`cat cozy-admin-passphrase`" cozy \
      cozy-stack instances add \
        --apps collect,contacts,drive,home,photos,settings,store,onboarding \
        --email <your email> example.com
    ```
   It's important to bear in mind here that Cozy instances are per-user: Each
   user requires their own subdomain and each instance will need to be created
   separately with this command.

5. Install the *onboarding* app with the following command:
    ```sh
    docker-compose exec -e "COZY_ADMIN_PASSWORD=`cat cozy-admin-passphrase`" cozy \
      cozy-stack apps install \
        --domain <example.com> \
        onboarding git://github.com/cozy/cozy-onboarding-v3.git#latest
    ```

6. The instance must now be registered. The previous command output a
   registration token. Navigate in your web browser to the following URL,
   substituting your domain and registration key at the appropriate locations:
     `https://<example.com>/?registerToken=<your token>`
   This will redirect you to a page where you will choose a password, and then
   you will be able to access your Cozy stack instance.

#### NOTE:
If you don't have access to the output of the `cozy-stack instances add`
command, you can get the registration token from `cozy-stack instances show <example.com>`.
The problem is, this version of the token is base-64 encoded, while we need the
hex-encoding. To convert the token to the appropriate format, use the following
command: `echo -n "<base64 token>" | base64 -d | xxd -p`
