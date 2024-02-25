# Using References

When working on a blueprint, we may often need to pull configuration data from multiple sources or navigate paths repeatedly.  
This can lead to duplication, making our blueprint harder to read and maintain, while increasing the risk of errors or inconsistencies when changes are made.

## So, what exactly are references? 

They are shortcuts to data defined elsewhere in the blueprint. References are particularly useful when the same data needs to be accessed in multiple places, reducing redundancy, improving readability, and making maintenance easier.
 
Let's dive into an example. Here's a typical blueprint without references:
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
 
The paths `$prodEnv.'configuration'.'api'` and `$prodEnv.'configuration'.'features'` are repeated multiple times.  
Instead of duplicating these, we can define references that point to them and reuse those references throughout the blueprint.  
 
We use the `$refs` attribute to do that. It's a dedicated section within the blueprint where we can define shortcuts to avoid repeating long and complex paths throughout our blueprint. Think of `$refs` as a container that holds all the references we want to use elsewhere in the blueprint.  
By organizing them all in one place, we make the blueprint cleaner and easier to comprehend.
```yaml
version: 1
sources:
  prodEnv:
    json:
      config:
        path: ./config.json
$refs:
```

First, let's create a reference to `$prodEnv.'configuration'.'api'`, and name it `apiConfig`:
```yaml
$refs:
  apiConfig: $prodEnv.'configuration'.'api'
``` 

A reference is a key-value pair, where the key is the name of it, and the value is a path to a node in the data we intend to use.  
Names should be meaningful, reflecting the data they represent, and can consist of letters and digits only.  

Let's add another, pointing to the `$prodEnv.'configuration'.'features'` node:
```yaml
$refs:
  apiConfig: $prodEnv.'configuration'.'api'
  featuresConfig: $prodEnv.'configuration'.'features'
``` 

Finally, let's add a reference to `$prodEnv.'configuration'.'features'.'themes'`:
```yaml
$refs:
  apiConfig: $prodEnv.'configuration'.'api'
  featuresConfig: $prodEnv.'configuration'.'features'
  themesFeatureConfig: $featuresConfig.'themes'
```  

<img alt="pointing fingers" src="https://wp.inews.co.uk/wp-content/uploads/2021/05/SEI_80560343.jpg?strip=all&quality=90" height="180px"/>  

That's right, we can also have a reference point to another.  
There's no need to bother about the order in which they're defined—**Conformize** will properly determine the order of evaluation.  
We just need to make sure there are no circular dependencies, if such occur, we will get an error message.

Now that we've defined our references, let's put them into action to see how they simplify our blueprint:
```yaml
ruleset:
  - $value: $apiConfig.'endpoint'
    equal:
      - value: https://api.yourapp.com
  - $value: $apiConfig.'timeout'
    withinRange:
      - value: 5000
      - value: 30000
  - $value: $featuresConfig.'auth'.'enabled'
    isTrue:
  - $value: $themesFeatureConfig.'enabled'
    isTrue:
  - $value: $themesFeatureConfig.'available'
    containsAllOf:
      - value:
          - light
          - dark
          - blue 
``` 
Instead of typing the full path each time, we use `$apiConfig`,`$featuresConfig` and `$themesFeatureConfig` as the root of our paths.

```yaml 
version: 1
sources:
  prodEnv:
    json:
      config:
        path: ./config.json

$refs:
  apiConfig: $prodEnv.'configuration'.'api'
  featuresConfig: $prodEnv.'configuration'.'features'
  themesFeatureConfig: $featuresConfig.'themes'
ruleset:
  - $value: $apiConfig.'endpoint'
    equal:
      - value: https://api.yourapp.com
  - $value: $apiConfig.'timeout'
    withinRange:
      - value: 5000
      - value: 30000
  - $value: $featuresConfig.'auth'.'enabled'
    isTrue:
  - $value: $themesFeatureConfig.'enabled'
    isTrue:
  - $value: $themesFeatureConfig.'available'
    containsAllOf:
      - value:
          - light
          - dark
          - blue 
``` 

### Why Use References?

**Clarity:** References make a blueprint concise and way easier to read.  
**Reusability:** References allow us to reuse the same paths across different rules, eliminating repetitive typing and ensuring consistency throughout the blueprint.  
**Maintainability:** If something changes in the configuration, we only need to tweak the reference in one place. Easy peasy!
 
While references definitely spiced up our food, we still felt like something is missing.
 
<img alt="Needs more salt!" src="https://media.makeameme.org/created/needs-more-salt-e9d22a74cf.jpg" height="160px">  
 
That's why we've introduced [scaffolding for blueprints](./scaffolding_a_blueprint.md).
