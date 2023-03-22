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
  
  postCustom(customData?: models.CustomData, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  postUpload(packageFiles?: models.RepoPackageFiles, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
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
    
  }

  namespace Models {
    
    type CustomData = {
  commit_sha: string;
  data: JSONObject;
  repo_name: string;
};
    
    type ErrorCode = ("InvalidID");
    
    type JSONObject = {
  [key: string]: {
  [key: string]: any;
};
};
    
    type RepoPackageFile = {
  error?: string;
  go_version?: string;
  name: string;
  packages: RepoPackages;
  path: string;
};
    
    type RepoPackageFiles = RepoPackageFile[];
    
    type RepoPackages = {
  dependencies?: string[];
  is_local?: boolean;
  name?: string;
  version?: string;
};
    
    type UnknownResponse = {
  body?: string;
  statusCode?: number;
};
    
  }
}

export = Breakdown;
