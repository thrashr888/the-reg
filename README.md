# Are you on The Reg?

A global service registry. Free public forwarding. $6.99/mo for unlimited private.

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
    redis-cli -h redis-my-first-redis.the-reg.name -p 6379
    $ reg name cf3f7336b1e0 cool-website
    http://cool-website.the-reg.name
    $ curl http://cool-website.the-reg.name
    <p>The server works!</p>
    $ reg name d39dd625947b cool-website
    https://cool-website.the-reg.name
    $ curl https://cool-website.the-reg.name
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
    $ reg get redis
    c65e2d0eb499.the-reg.name:6379
    # redis.thrashr888.the-reg.name:6379
```

# Ideas

- Use Consul for secure routing/networking/proxying?
- Use ngrok for opening tunnels?
