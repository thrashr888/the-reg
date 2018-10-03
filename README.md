# Are you on The Reg?

A global service registry. Free public forwarding. $6.99/mo for unlimited private.

# CLI Help

```
$ reg [command]

account - `reg account new :username :email` sign up for an account
add - `reg add :name [hostname] :port` add a node
create - get a user token
get - `reg get :name` Get a service url
help - show this list
ip - get your public ip address
list - list your nodes
login - save your auth token
me - your username
name - `reg name :id :name` name a node
start - attempt to reset status to "UP"
server - run the web service
```

# API

```
GET /node
POST /node
GET /node/:id
PATCH /node/:id
DELETE /node/:id

GET /account
GET /account/confirm/:token

POST /user
GET /user/:id
PATCH /user/:id
DELETE /user/:id
```

# CLI Examples

```shell
    $ reg create
    # echo "authtoken: Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7" > ~/.thereg.yml
    # export THE_REG_TOKEN=Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7
    $ reg login Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7
    # echo "authtoken: Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7" > ~/.thereg.yml
    # export THE_REG_TOKEN=Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7
    $ reg me
    full-buffallo-hotness
    $ reg ip
    76.87.249.25
    $ reg add redis 6379 --public
    c65e2d0eb499
    $ reg list
    ID             NAME               HOST          PORT      STATUS   AGE   PUBLIC   TAGS
    dff8522fe5dc   first-deployment   76.87.249.25  8080:80   UP       6h    N        
    cf3f7336b1e0   http               76.87.249.25  80        UP       3h    Y        
    d39dd625947b   https              76.87.249.25  443       UP       3h    Y        
    bc2740d30a5f   httpexposed        76.87.249.25  8081:80   DOWN     2h    Y        
    c65e2d0eb499   redis              76.87.249.25  6379      UP       2h    Y        
    $ reg start httpexposed
    Local port 8081 not found. Try restarting your server.
    $ reg get redis
    c65e2d0eb499.the-reg.name:6379
    # redis.full-buffallo-hotness.the-reg.name:6379
    $ reg get first-deployment
    http://dff8522fe5dc.the-reg.name
    # http://first-deployment.full-buffallo-hotness.the-reg.name
    $ reg name c65e2d0eb499 redis-my-first-redis
    redis-cli -h redis-my-first-redis.full-buffallo-hotness.the-reg.name -p 6379
    $ reg name cf3f7336b1e0 cool-website
    http://cool-website.full-buffallo-hotness.the-reg.name
    $ curl http://cool-website.full-buffallo-hotness.the-reg.name
    <p>The server works!</p>
    $ reg name d39dd625947b cool-website
    https://cool-website.full-buffallo-hotness.the-reg.name
    $ curl https://cool-website.full-buffallo-hotness.the-reg.name
    <p>The server works!</p>
    $ reg add redis-b redis15.localnet.org 6390 --public
    y0am0fa6786a
    $ reg list
    ID             NAME               HOST                  PORT      STATUS   AGE   PUBLIC   TAGS
    dff8522fe5dc   first-deployment   76.87.249.25          8080:80   UP       6h    N        
    cf3f7336b1e0   http               76.87.249.25          80        UP       3h    Y        
    d39dd625947b   https              76.87.249.25          443       UP       3h    Y        
    bc2740d30a5f   httpexposed        76.87.249.25          8081:80   DOWN     2h    Y        
    c65e2d0eb499   redis              76.87.249.25          6379      UP       2h    Y        
    y0am0fa6786a   redis-b            redis15.localnet.org  6379      UP       2h    Y        
    $ reg account new thrashr888 thrashr888@gmail.com
    Account created. Check your email to log in at https://www.the-reg.name/
    # You click the link to confirm your email address...
    $ reg get redis
    c65e2d0eb499.full-buffallo-hotness.the-reg.name:6379
    # redis.thrashr888.the-reg.name:6379
```

# Ideas

- Use Consul for secure routing/networking/proxying?
- Use ngrok for opening tunnels?

# FAQ

**Does this use UDDI?**

No way.

**Is this a totally secure service?**

It's better if you assume it is not.

**Can people use this service to log into my computer/server?**

It's more like some links or a [proxy server](https://en.wikipedia.org/wiki/Proxy_server). It's up to you to firewall or otherwise secure your servers.

# Dev

> reflex -r '\.(go|html)$' -s -- sh -c 'go run reg.go'
