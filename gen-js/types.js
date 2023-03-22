module.exports.Errors = {};

/**
 * BadRequest
 * @extends Error
 * @memberof module:breakdown
 * @alias module:breakdown.Errors.BadRequest
 * @property  code
 * @property {string} message
 */
module.exports.Errors.BadRequest = class extends Error {
  constructor(body) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

/**
 * InternalError
 * @extends Error
 * @memberof module:breakdown
 * @alias module:breakdown.Errors.InternalError
 * @property  code
 * @property {string} message
 */
module.exports.Errors.InternalError = class extends Error {
  constructor(body) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

