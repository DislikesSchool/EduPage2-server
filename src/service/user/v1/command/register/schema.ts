import { extendApi } from '@purista/core';
import { z } from 'zod';

export const userV1RegisterInputParameterSchema = extendApi(z.any(), {
  title: 'User registration parameters',
});

// define the input payload
export const userV1RegisterInputPayloadSchema = extendApi(
  z.object({
    username: extendApi(z.string().nonempty(), { title: 'Username' }),
    password: extendApi(z.string().nonempty(), { title: 'Password' }),
  }),
  {
    title: 'User registration payload',
  },
);

// define the output payload
export const userV1RegisterOutputPayloadSchema = extendApi(
  z.object({
    token: extendApi(z.string().nonempty(), { title: 'User token' }),
  }),
  {
    title: 'User registration response',
  },
);
