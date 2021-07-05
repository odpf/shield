import Confidence from 'confidence';
import Inert from '@hapi/inert';
import Vision from '@hapi/vision';
import Glue from '@hapi/glue';
import H2o2 from '@hapi/h2o2';
import HapiSwagger from 'hapi-swagger';
import Qs from 'qs';
import * as Config from './config';
import Logging from '../plugin/logging';

const internals: any = {
  criteria: {
    env: process.env.NODE_ENV
  }
};

internals.manifest = {
  $meta: 'App manifest',
  server: {
    port: Config.get('/port/web'),
    query: {
      parser: (query: string) =>
        Qs.parse(query, { comma: true, allowDots: true })
    },
    router: {
      stripTrailingSlash: true,
      isCaseSensitive: false
    },
    compression: false,
    routes: {
      security: true,
      cors: true
    }
  },
  register: {
    plugins: [
      {
        plugin: '../plugin/postgres'
      },
      {
        plugin: H2o2
      },
      {
        plugin: '../plugin/iap'
      },
      {
        plugin: '../plugin/events'
      },
      {
        plugin: '../plugin/iam/index',
        options: Config.get('/postgres')
      },
      {
        plugin: '../app/ping/index'
      },
      {
        plugin: '../app/profile/index'
      },
      {
        plugin: '../app/user/index'
      },
      {
        plugin: '../app/group/index'
      },
      {
        plugin: '../app/role/index'
      },
      {
        plugin: '../app/activity/index'
      },
      {
        plugin: '../app/access/index'
      },
      {
        plugin: '../app/resource/index'
      },
      {
        plugin: '../plugin/proxy'
      },
      {
        plugin: Inert
      },
      {
        plugin: Vision
      },
      {
        plugin: HapiSwagger,
        options: {
          info: {
            title: 'API',
            version: '0.1.0'
          },
          securityDefinitions: {
            AUTH_KEY: {
              type: 'apiKey', // apiKey is defined by the Swagger spec
              name: 'X-Goog-Authenticated-User-Email', // the name of the query parameter / header
              in: 'header'
            }
          },
          security: [{ AUTH_KEY: [] }],
          grouping: 'tags'
        }
      },
      { ...Logging }
    ]
  }
};

internals.store = new Confidence.Store(internals.manifest);

const serverManifest = internals.store.get('/', internals.criteria);

const options = {
  relativeTo: __dirname
};

export default Glue.compose.bind(Glue, serverManifest, options);
