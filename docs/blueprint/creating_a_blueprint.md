# Creating a Blueprint

All we need is a trusty text editor. While it's handy if it supports YAML or JSON syntax highlighting, any editor will do.  
There aren't any plugins currently available to assist with blueprints, but that's okay—we'll walk through everything step-by-step. 

Once our editor is ready, it's time to dive in. This guide will cover everything we need to create a blueprint, ensuring we can confidently validate our configurations. So, let's buckle up and get started!  

Before we create a blueprint, let's take a closer look at the JSON configuration file we'll be working with. Imagine this file is for a web application managing user authentication, API connections, and theme settings. No need to imagine—we can download the sample file from [this link](../samples/config.json) and save it to a preferred location.
  
The file includes settings for API endpoints, user authentication methods, available themes, and database connections.  
Our goal? To ensure these configurations are rock solid before releasing our application into the wild.

First things first, we need to specify the version of our blueprint schema. In our editor, we add the following:
```yaml
version: 1
```

Now that we've laid the foundation, let's move on to the next step.  
Next, we need to define in our blueprint where to pull the configuration data from—in this case, our `config.json` file.  
We do this by adding the following lines:
```yaml
sources:
  prodEnv:
    json:
      config:
        path: ./config.json
```

Here, `prodEnv` is an alias for the production environment configuration.  
The provider is set to `json`, and `path` points to the location of our configuration file. 
  
Now, we define the ruleset, telling **Conformize** what checks to perform on the data in our JSON file.  
Let's start with a rule to ensure the API endpoint is set correctly:
```yaml
ruleset:
  - $value: $prodEnv.'configuration'.'api'.'endpoint'
    equal:
      - value: https://api.yourapp.com
```

Next, we'll ensure the API timeout is within a safe range—between 5,000 and 30,000 milliseconds:
```yaml
  - $value: $prodEnv.'configuration'.'api'.'timeout'
    withinRange:
      - value: 5000
      - value: 30000
```

We also want to confirm that authentication is enabled before going live:
```yaml
  - $value: $prodEnv.'configuration'.'features'.'auth'.'enabled'
    isTrue:
```
 
Finally, let's ensure our newly added themes feature is not only enabled but also includes a list of available themes:
```yaml
  - $value: $prodEnv.'configuration'.'features'.'themes'.'enabled'
    isTrue:
  - $value: $prodEnv.'configuration'.'features'.'themes'.'available'
    containsAllOf:
      - value:
          - light
          - dark
          - blue
```  
 
## Putting It All Together  

After defining the version, sources, and rules, our blueprint should look like this:
```yaml
version: 1
sources:
  prodEnv:
    json:
      config:
        path: ./config.json
ruleset:
  - $value: $prodEnv.'configuration'.'api'.'endpoint'
    equal:
      - value: https://api.yourapp.com
  - $value: $prodEnv.'configuration'.'api'.'timeout'
    withinRange:
      - value: 5000
      - value: 30000
  - $value: $prodEnv.'configuration'.'features'.'auth'.'enabled'
    isTrue:
  - $value: $prodEnv.'configuration'.'features'.'themes'.'enabled'
    isTrue:
  - $value: $prodEnv.'configuration'.'features'.'themes'.'available'
    containsAllOf:
      - value:
          - light
          - dark
          - blue
``` 

Save it to a file named **blueprint.cnfrm.yaml**, in the same location as the sample configuration file.

With our blueprint ready, it's time to validate our configuration using **Conformize**. Open a terminal, navigate to the directory where the blueprint is saved, and run:
```
$ conformize blueprint apply
```

If everything is set up correctly, we should see output like this shorty after running the command:
```
Configuring source 'prodEnv' with 'json' provider...
Configured source 'prodEnv' with 'json' provider.

Reading from source 'prodEnv'...
Done reading from source 'prodEnv'.

Blueprint has been applied successfully.
```

The message **`Blueprint has been applied successfully.`** confirms that the configuration matches the expectations set out in our blueprint.

That's a rather simple example, in reality our ruleset could contain many more rules, and one thing becomes apparent—it's irritating having to write repeating paths again and again:   
`$prodEnv.'configuration'.'api'`.'endpoint'  
`$prodEnv.'configuration'.'api'`.'timeout'  
`$prodEnv.'configuration'.'features'`.'auth'.'enabled'  
`$prodEnv.'configuration'.'features'`.'themes'.'enabled'  
`$prodEnv.'configuration'.'features'`.'themes'.'available'  

<img alt="if we could avoid writing the full path over and over, that woudl be great" src="https://i.imgflip.com/91znvv.jpg" height="200px"/>

We eat our own dog food here, and we didn't like the taste of it—so we were eager to introduce a way to ease working with paths.  
It's about time we become familiar with another attribute of a blueprint - [references](./using_references.md).