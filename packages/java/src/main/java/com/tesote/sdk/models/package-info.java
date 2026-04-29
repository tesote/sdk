/**
 * Typed model records for every API resource. Records are immutable, mirror the
 * API JSON one-to-one (snake_case wire names mapped to camelCase Java fields
 * via Jackson's {@code @JsonProperty}), and tolerate unknown fields so the
 * API may add new keys without breaking older SDK builds.
 */
package com.tesote.sdk.models;
