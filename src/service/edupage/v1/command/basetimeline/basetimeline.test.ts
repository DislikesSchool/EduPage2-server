import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { edupageV1Service } from '../../edupageV1Service';
import { basetimelineCommandBuilder } from './basetimelineCommandBuilder';
import {
  EdupageV1BasetimelineInputParameter,
  EdupageV1BasetimelineInputPayload,
} from './types';

describe('service Edupage version 1 - command basetimeline', () => {
  let sandbox = createSandbox();
  beforeEach(() => {
    sandbox = createSandbox();
  });

  afterEach(() => {
    sandbox.restore();
  });

  test('does not throw', async () => {
    const service = edupageV1Service.getInstance(
      getEventBridgeMock(sandbox).mock,
      { logger: getLoggerMock(sandbox).mock },
    );

    const basetimeline = basetimelineCommandBuilder
      .getCommandFunction()
      .bind(service);

    const payload: EdupageV1BasetimelineInputPayload = undefined;

    const parameter: EdupageV1BasetimelineInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await basetimeline(context.mock, payload, parameter);

    expect(result).toBeUndefined();
  });
});
