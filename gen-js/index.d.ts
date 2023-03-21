import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
  timeout?: number;
  retryPolicy?: RetryPolicy;
}

interface IterResult<R> {
  map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
  toArray(cb?: Callback<R[]>): Promise<R[]>;
  forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
  forEachAsync(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
}

interface CircuitOptions {
  forceClosed?: boolean;
  maxConcurrentRequests?: number;
  requestVolumeThreshold?: number;
  sleepWindow?: number;
  errorPercentThreshold?: number;
}

interface GenericOptions {
  timeout?: number;
  keepalive?: boolean;
  retryPolicy?: RetryPolicy;
  logger?: Logger;
  circuit?: CircuitOptions;
  serviceName?: string;
}

interface DiscoveryOptions {
  discovery: true;
  address?: undefined;
}

interface AddressOptions {
  discovery?: false;
  address: string;
}

type BreakdownOptions = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = Breakdown.Models

declare class Breakdown {
  constructor(options: BreakdownOptions);

  close(): void;
  
  healthCheck(options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  getThings(options?: RequestOptions, cb?: Callback<models.Thing[]>): Promise<models.Thing[]>
  
  deleteThing(id: string, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  getThing(id: string, options?: RequestOptions, cb?: Callback<models.Thing>): Promise<models.Thing>
  
  createOrUpdateThing(params: models.CreateOrUpdateThingParams, options?: RequestOptions, cb?: Callback<models.Thing>): Promise<models.Thing>
  
}

declare namespace Breakdown {
  const RetryPolicies: {
    Single: RetryPolicy;
    Exponential: RetryPolicy;
    None: RetryPolicy;
  }

  const DefaultCircuitOptions: CircuitOptions;

  namespace Errors {
    interface ErrorBody {
      message: string;
      [key: string]: any;
    }

    
    class BadRequest {
  code?: models.ErrorCode;
  message?: string;

  constructor(body: ErrorBody);
}
    
    class InternalError {
  code?: models.ErrorCode;
  message?: string;

  constructor(body: ErrorBody);
}
    
    class NotFound {
  code?: models.ErrorCode;
  message?: string;

  constructor(body: ErrorBody);
}
    
  }

  namespace Models {
    
    type CreateOrUpdateThingParams = {
  thing?: Thing;
  id: string;
};
    
    type ErrorCode = ("InvalidID");
    
    type Thing = {
  foo?: string;
  id?: string;
};
    
    type UnknownResponse = {
  body?: string;
  statusCode?: number;
};
    
  }
}

export = Breakdown;
