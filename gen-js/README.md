<a name="module_breakdown"></a>

## breakdown
breakdown client library.


* [breakdown](#module_breakdown)
    * [Breakdown](#exp_module_breakdown--Breakdown) ⏏
        * [new Breakdown(options)](#new_module_breakdown--Breakdown_new)
        * _instance_
            * [.close()](#module_breakdown--Breakdown+close)
            * [.healthCheck([options], [cb])](#module_breakdown--Breakdown+healthCheck) ⇒ <code>Promise</code>
            * [.postCustom(customData, [options], [cb])](#module_breakdown--Breakdown+postCustom) ⇒ <code>Promise</code>
            * [.postDeploy(deploys, [options], [cb])](#module_breakdown--Breakdown+postDeploy) ⇒ <code>Promise</code>
            * [.postUpload(repoCommit, [options], [cb])](#module_breakdown--Breakdown+postUpload) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_breakdown--Breakdown.RetryPolicies)
                * [.Exponential](#module_breakdown--Breakdown.RetryPolicies.Exponential)
                * [.Single](#module_breakdown--Breakdown.RetryPolicies.Single)
                * [.None](#module_breakdown--Breakdown.RetryPolicies.None)
            * [.Errors](#module_breakdown--Breakdown.Errors)
                * [.BadRequest](#module_breakdown--Breakdown.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_breakdown--Breakdown.Errors.InternalError) ⇐ <code>Error</code>
            * [.DefaultCircuitOptions](#module_breakdown--Breakdown.DefaultCircuitOptions)

<a name="exp_module_breakdown--Breakdown"></a>

### Breakdown ⏏
breakdown client

**Kind**: Exported class  
<a name="new_module_breakdown--Breakdown_new"></a>

#### new Breakdown(options)
Create a new client object.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| options | <code>Object</code> |  | Options for constructing a client object. |
| [options.address] | <code>string</code> |  | URL where the server is located. Must provide this or the discovery argument |
| [options.discovery] | <code>bool</code> |  | Use clever-discovery to locate the server. Must provide this or the address argument |
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. |
| [options.keepalive] | <code>bool</code> |  | Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies) | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| [options.logger] | <code>module:kayvee.Logger</code> | <code>logger.New(&quot;breakdown-wagclient&quot;)</code> | The Kayvee logger to use in the client. |
| [options.circuit] | <code>Object</code> |  | Options for constructing the client's circuit breaker. |
| [options.circuit.forceClosed] | <code>bool</code> |  | When set to true the circuit will always be closed. Default: true. |
| [options.circuit.maxConcurrentRequests] | <code>number</code> |  | the maximum number of concurrent requests the client can make at the same time. Default: 100. |
| [options.circuit.requestVolumeThreshold] | <code>number</code> |  | The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. |
| [options.circuit.sleepWindow] | <code>number</code> |  | how long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. |
| [options.circuit.errorPercentThreshold] | <code>number</code> |  | the threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. |

<a name="module_breakdown--Breakdown+close"></a>

#### breakdown.close()
Releases handles used in client

**Kind**: instance method of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  
<a name="module_breakdown--Breakdown+healthCheck"></a>

#### breakdown.healthCheck([options], [cb]) ⇒ <code>Promise</code>
Checks if the service is healthy

**Kind**: instance method of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_breakdown--Breakdown.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_breakdown--Breakdown.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_breakdown--Breakdown+postCustom"></a>

#### breakdown.postCustom(customData, [options], [cb]) ⇒ <code>Promise</code>
upload or replace custom data for a given repo and commit SHA

**Kind**: instance method of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_breakdown--Breakdown.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_breakdown--Breakdown.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| customData |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_breakdown--Breakdown+postDeploy"></a>

#### breakdown.postDeploy(deploys, [options], [cb]) ⇒ <code>Promise</code>
report a number of deploys

**Kind**: instance method of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_breakdown--Breakdown.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_breakdown--Breakdown.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| deploys |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_breakdown--Breakdown+postUpload"></a>

#### breakdown.postUpload(repoCommit, [options], [cb]) ⇒ <code>Promise</code>
upload a package-type file, generated by breakdown-cli

**Kind**: instance method of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_breakdown--Breakdown.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_breakdown--Breakdown.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| repoCommit |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_breakdown--Breakdown.RetryPolicies"></a>

#### Breakdown.RetryPolicies
Retry policies available to use.

**Kind**: static property of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  

* [.RetryPolicies](#module_breakdown--Breakdown.RetryPolicies)
    * [.Exponential](#module_breakdown--Breakdown.RetryPolicies.Exponential)
    * [.Single](#module_breakdown--Breakdown.RetryPolicies.Single)
    * [.None](#module_breakdown--Breakdown.RetryPolicies.None)

<a name="module_breakdown--Breakdown.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies)  
<a name="module_breakdown--Breakdown.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies)  
<a name="module_breakdown--Breakdown.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_breakdown--Breakdown.RetryPolicies)  
<a name="module_breakdown--Breakdown.Errors"></a>

#### Breakdown.Errors
Errors returned by methods.

**Kind**: static property of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  

* [.Errors](#module_breakdown--Breakdown.Errors)
    * [.BadRequest](#module_breakdown--Breakdown.Errors.BadRequest) ⇐ <code>Error</code>
    * [.InternalError](#module_breakdown--Breakdown.Errors.InternalError) ⇐ <code>Error</code>

<a name="module_breakdown--Breakdown.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

**Kind**: static class of [<code>Errors</code>](#module_breakdown--Breakdown.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code |  | 
| message | <code>string</code> | 

<a name="module_breakdown--Breakdown.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of [<code>Errors</code>](#module_breakdown--Breakdown.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code |  | 
| message | <code>string</code> | 

<a name="module_breakdown--Breakdown.DefaultCircuitOptions"></a>

#### Breakdown.DefaultCircuitOptions
Default circuit breaker options.

**Kind**: static constant of [<code>Breakdown</code>](#exp_module_breakdown--Breakdown)  
