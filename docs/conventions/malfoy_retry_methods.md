Here i'll document the retry methods used in malfoy chamber.

## none

Default method, no retry will be done.

## fixed

Retry will be done after a fixed delay, defined by the user in the `RETRY_DELAY` option.

## linear

Retry will be done after a delay that increases linearly with each attempt. The delay is calculated as `RETRY_DELAY * attempt_number`, where `attempt_number` starts at 1 for the first retry.

## exponential

Retry will be done after a delay that increases exponentially. The delay is calculated as `RETRY_DELAY * (2 ^ attempt_number)`, where `attempt_number` starts at 1 for the first retry.

## exponential_jitter

Similar to exponential, but with added randomness to avoid thundering herd problems. The delay is calculated as `RETRY_DELAY * ((2 ^ attempt_number) ± random_jitter)`, where `random_jitter` is set by the user in the `JITTER` option as a percentage of the calculated delay.

