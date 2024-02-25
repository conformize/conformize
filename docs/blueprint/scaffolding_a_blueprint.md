# Scaffolding a Blueprint

If only there was an easier way to create a blueprint... Well, there is!  

Scaffolding in Conformize allows us to quickly create base structure for a blueprint, saving time and ensuring consistency. 
Instead of starting from scratch, we can use the scaffold command to create a template that we can then fill in with specific details.

To learn more about the scaffold command, we must open a terminal and run:

```
$ ./conformize scaffold blueprint help
```

A brief description of the command and the available arguments that can be used is presented to us:

```
description: create a blueprint scaffold

usage: conformize [global options] scaffold blueprint [subcommand] [arguments]

available subcommands:

help                  display usage instructions

available arguments:

-format               specifies the output format for blueprint scaffold - JSON or YAML, e.g. -format yaml. YAML format will be used if not specified.
-predicates           a comma-separated list of predicates to ba added to the blueprint scaffold
-provider             the provider to be used to retrieve data from a configuration source
-refs                 a comma-separated list of reference aliases to be defined in the blueprint scaffold
-source               the alias for a configuration source to be added to blueprint scaffold
-version              specifies the schema version to be used for blueprint scaffold. default value of 1 will be used unless provided.
```

Let's break down a few key arguments:

* **-source**: Specifies the configuration source alias to be scaffolded. For every source argument, we must pair it with a -provider argument.

* **-provider**: Specifies the backend (e.g., JSON, YAML, or AWS Parameter Store) that will be used to retrieve data from the configuration source.

* **-predicates**: Specifies the conditions or rules that will be applied to configuration values.

* **-refs**: Specifies the reference aliases that will be added to the scaffold.

* **-format**: Specifies whether the blueprint should be generated in JSON or YAML format. YAML is the default if this argument is not provided.

Let's say we want to create a blueprint for configuration data scattered across multiple locations, such as a YAML file and AWS SSM Parameter Store. Here's how we can scaffold that:

```
$ conformize scaffold blueprint \
    -source uatEnvAwsParamStore -provider aws_parameter_store \
    -source uatEnvAppConfigFile -provider yaml
```

We are presented with the following output:

```
Blueprint scaffold created at ./blueprint.cnfrm.yaml.
```

When we open the newly created `blueprint.cnfrm.yaml` file, we'll see the following content:

```yaml
version: 1
sources:
  uatEnvAppConfigFile:
    yaml:
  uatEnvAwsParamStore:
    aws_parameter_store:
```  
 
A few things to note:

* Since no version was specified, version: 1 was added by default.
* The sources section lists each configuration source and its associated provider.

To retrieve data from these sources, we'll need to add [configuration](../providers/configuring_a_provider.md) for the specified providers. 
 
> Important Note: The scaffold will be re-created every time we run the scaffold command. Any manual changes made to the file after it's generated will be overwritten. We must make sure to save our changes or move the file before re-running the scaffold command. 

Next, let's add some references to our blueprint scaffold. Here's an example:

```
$ conformize scaffold blueprint \
    -source uatEnvAwsParamStore -provider aws_parameter_store \
    -source uatEnvAppConfigFile -provider yaml \
    -refs uatEnvAppFeaturesConfig,uatEnvAppApiConfig
```

After running the command, the updated `blueprint.cnfrm.yaml` file should look like this:

```yaml
version: 1
sources:
  uatEnvAppConfigFile:
    yaml:
  uatEnvAwsParamStore:
    aws_parameter_store:
$refs:
  uatEnvAppApiConfig:
  uatEnvAppFeaturesConfig:
```

The `$refs` section contains scaffolded references for each alias provided.  
We'll need to manually fill in the actual path to the data each reference should point to..

Let's now look at the blueprint example below.

```yaml
version: 1
sources:
  prodEnv:
    json:
      config:
        path: ./config.json

$refs:
  prodApiConfig: $prodEnv.'configuration'.'api'
  prodFeaturesConfig: $prodEnv.'configuration'.'features'

ruleset:
  - $value: $prodApiConfig.'endpoint'
    equal:
      - value: https://api.yourapp.com
  - $value: $prodApiConfig.'timeout'
    withinRange:
      - value: 5000
      - value: 30000
  - $value: $prodFeaturesConfig.'auth'.'enabled'
    isTrue:
  - $value: $prodFeaturesConfig.'themes'.'enabled'
    isTrue:
  - $value: $prodFeaturesConfig.'themes'.'available'
    containsAllOf:
      - value:
          - light
          - dark
          - blue 
```

We can create a scaffold for this blueprint by running the following command:

```
$ conformize scaffold blueprint \
  -source stageEnv -provider consul \
  -source prodEnv -provider consul \
  -refs stageApiConfig,prodApiConfig \
  -predicates equal,withinRange,isTrue,isTrue,containsAllOf
```

Once created, our `blueprint.cnfrm.yaml` looks like this:

```yaml
version: 1
sources:
  prodEnv:
    consul:
      config:
  stageEnv:
    consul:
      config:
$refs:
  prodApiConfig:
  stageApiConfig:
ruleset:
  - $value:
    equal:
      - value:
  - $value:
    withinRange:
      - value:
      - value:
  - $value:
    isTrue:
  - $value:
    isTrue:
  - $value:
    containsAllOf:
      - value:
```  

In this scaffold we see the following: 

* A ruleset section that defines validation checks based on the predicates we provided (`equal`, `withinRange`, etc.).
* Each rule needs to have `$value` filled-in with the actual path to the data we want to validate.
* Whenever the predicate requires arguments (e.g. `withinRange`), we'll see added arguments placeholders with blank values, which should be replaced with the actual ones.

<img alt="it ain't much" src="https://media.npr.org/assets/img/2023/05/26/honest-work-meme-c7034f8bd7b11467e1bfbe14b87a5f6a14a5274b.jpg" height="180px">  

And while scaffolding gives us a head start, we still need to make manual edits to complete the blueprint. 
As **Conformize** evolves, scaffolding will become more powerful, allowing for more automation and flexibility in the future.
 
## Why Use Scaffolding? 

* **Saves Time:** Scaffolding provides a base blueprint structure, so we don't have to build it from scratch.
* **Consistency:** It ensures that our blueprint follows a standardized format, reducing the chances of errors.
* **Extensibility:** We can easily add multiple sources, references, and rules by specifying command arguments.
