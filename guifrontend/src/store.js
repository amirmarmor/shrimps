import { configureStore, getDefaultMiddleware } from '@reduxjs/toolkit'
import configReducer from './features/config/configSlice'

export default function configureAppStore(preloadedState) {
  return configureStore({
    reducer: {
      config: configReducer
    },
    middleware: [...getDefaultMiddleware()],
    preloadedState,
  })

  // if (process.env.NODE_ENV !== 'production' && module.hot) {
  //   module.hot.accept('./features/config/', () => store.replaceReducer(rootReducer))
  // }
}