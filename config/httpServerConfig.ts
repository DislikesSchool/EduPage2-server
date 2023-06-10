import { HttpServerServiceV1Config } from '@purista/httpserver';

const httpServerConfig: HttpServerServiceV1Config = {
  fastify: {
    /* ssl config
    http2: true,
    https: {
      allowHTTP1: true,
      ca: ``, // add your certificate content here
      cert: ``, // add your certificate content here
    },
    */
  },
  port: 80,
  logLevel: 'debug',
  domain: 'localhost',
  host: '',
  cookieSecret: 'oCrUlLnZqhj99evenJ3x',
  apiMountPath: '/api',
  openApi: {
    enabled: true,
    info: {
      title: 'EduPage2 public API',
      description: 'OpenApi definition for EduPage2 public endpoints',
      version: '1.0.0',
    },
  },
};

export default httpServerConfig;
