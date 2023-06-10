import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { userV1Service } from '../../userV1Service';
import { registerCommandBuilder } from './registerCommandBuilder';
import {
  UserV1RegisterInputParameter,
  UserV1RegisterInputPayload,
} from './types';

describe('service User version 1 - command register', () => {
  let sandbox = createSandbox();
  beforeEach(() => {
    sandbox = createSandbox();
  });

  afterEach(() => {
    sandbox.restore();
  });

  test('does not throw', async () => {
    const service = userV1Service.getInstance(
      getEventBridgeMock(sandbox).mock,
      { logger: getLoggerMock(sandbox).mock },
    );

    const register = registerCommandBuilder.getCommandFunction().bind(service);

    const payload: UserV1RegisterInputPayload = {
      username: 'a',
      password: 'b',
    };

    const parameter: UserV1RegisterInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await register(context.mock, payload, parameter);

    expect(result).toBeDefined();
  });
});
