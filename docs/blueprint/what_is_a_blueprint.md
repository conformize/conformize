# What do we mean by a Blueprint?

A **Blueprint** is our detailed game plan to ensure every configuration stays in check before the big release.  
Just like a building's blueprint is a detailed sketch that follows established standards, our sketch is a YAML or JSON file with specific attributes and values.

## Attributes

A **Blueprint** consists of several key attributes that define how **Conformize** should operate. 
Let's learn more about their meaning. 

>Provided examples use **YAML** but bear in mind that a Blueprint can be defined by using either **JSON** or **YAML** format.

### 1. version
- **Description:** Specifies the version of the Blueprint schema to be used.
- **Example:**   
`version: 1`
> Current and only supported version for the time being is **1**.

### 2. sources
- **Description:** Defines the configuration sources that data will be retrieved from.
Each key in it acts as an alias for a source that can be referenced later to query data.
- **Example:**
```yaml
sources:
  devEnv:
    json:
      config:
        path: ./config/environment.dev.json
```

As we can see in the example above, an alias is a dictionary with a single key specifying the provider to be used for retrieving data. 
The value of the `config` key will vary depending on the backend implementation. To learn more about how to configure and use providers, refer to the [built-in providers](./providers/built-in.md) documentation.

### 3. ruleset
- **Description:** The ruleset in a Blueprint is where the actual validation logic lives.  
It's a collection of rules that **Conformize** will apply to our configuration data. Each rule defines a specific check or condition that must be met, ensuring our configurations are valid and conform to expectations.
- **Example:**
```yaml
ruleset:
  - $value: $devEnv.'configuration'.'api'.'endpoint'
    equal:
      - value: https://api.example.com
  - $value: $devEnv.'configuration'.'api'.'timeout'
    greaterThanOrEqual:
      - value: 5000
  - $value: $devEnv.'configuration'.'features'.'auth'.'enabled'
    equal:
      - value: true
  - $value: $devEnv.'configuration'.'database'.'pool'.'min'
    withinRange:
      - value: 5
      - value: 10
  - $value: $devEnv.'configuration'.'features'.'themes'.'available'
    containsAllOf:
      - value:
          - light
          - dark
          - blue
```  

<img src="https://i.imgflip.com/91chlb.jpg" height="220px"/>

## Defining Rules and Predicates  
A rule consists of the following parts:

### **$value**
- **Description:** Specifies the path to the value that the rule should validate.  
We begin the path with the dollar sign `$`, which denotes the root element, followed by keys in single quotes that lead to the specific data point.

- **Example:**  
```yaml
$value: $devEnv.'configuration'.'api'.'timeout'
```  
  
Here, we target the `timeout` value in the API configuration of the `devEnv` source.

### **predicate**
- **Description:** A predicate defines the condition a value needs to meet—whether it's an equality check, a comparison, or another type of validation.
- **Example:**  
```yaml
$value: $devEnv.'configuration'.'api'.'endpoint'
equal:
  - value: https://api.example.com
``` 
In this example, the `equal` predicate checks if the value equals the specified argument.

When a predicate needs arguments, they're given directly as a list under it.
These arguments specify the value or path that the predicate will use to perform the validation.

Here are the types of arguments we can use to define a predicate:

#### **1. value**  
- **Description**: A raw value to compare against. 
- **Example**:
```yaml
greaterThanOrEqual:
  - value: 5000 
```

#### **2. path**  
- **Description**: A reference to another value within the data to compare with.
- **Example**:
```yaml
notEqual:
  - path: $stageEnv.'appConfig'.'api'.'endpoint'
```

#### **3. sensitive**  
- **Description**: We use it to hide sensitive values from logs or outputs.
- **Example**:
```yaml
notEqual:
  - sensitive:
      value: https://api.example.com
```

Now that we've nailed the details, let's go ahead and [create our first blueprint](./creating_a_blueprint.md)!