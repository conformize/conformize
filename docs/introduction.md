# Introduction

## Why do we build it?

Building software has become tremendously complex since the beginning of this century.  
Software design paradigms such as Service-Oriented Architecture (SOA), have emerged, adding complexity.    
Systems are becoming more and more sophisticated, having many interconnected components, be it in-house built or third-party ones.  
They often require various configurations, which might be scattered across multiple locations, making management far from simple.  

It's great when everything just works, but things change—and they can do so frequently.  
Development is complete, tests and reviews have passed, and we deploy the latest version of our product.  

However, instead of receiving the anticipated praise of users, complaints start pouring in and alerts to appear in our monitoring.  
This scenario is familiar to many of us, as we've often traced issues back to a simple misconfiguration. It's not just frustrating, but can also  
lead to costly consequences like downtime and loss of user trust.

## What are we up to, then?

Much like how a test framework helps us ensure code behaves as expected, **Conformize** lets us write scenarios to  
validate configuration structure and data.

## But how do we do it?!

The four pillar components of **Conformize** are the Configuration Data Provider, Type System, Predicate, and Blueprint.  
Below is a brief introduction to each:

* **Configuration Data Provider** - Fetches data and transforms it into a unified structure that other components can work with.
* **Type System** - With no predefined data schema required, we're diving into the unknown, figuring out what type of data  
we're dealing with and turning it into something our system can understand.
* **Predicate** - Ensures value plays by the rules and meets the conditions we've set.
* **Blueprint** - The cornerstone of the solution, enabling us to map out data sources, define queries, and lay down the rules to follow.

[Let's get it started!](./getting_started/getting_started.md)