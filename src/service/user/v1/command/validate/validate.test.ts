import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { userV1Service } from '../../userV1Service';
import {
  UserV1ValidateInputParameter,
  UserV1ValidateInputPayload,
} from './types';
import { validateCommandBuilder } from './validateCommandBuilder';

describe('service User version 1 - command validate', () => {
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

    const validate = validateCommandBuilder.getCommandFunction().bind(service);

    const payload: UserV1ValidateInputPayload = {
      token: 'abcd',
    };

    const parameter: UserV1ValidateInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await validate(context.mock, payload, parameter);

    expect(result).toBeDefined();
  });
});
