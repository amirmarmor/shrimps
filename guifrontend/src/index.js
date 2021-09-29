import React from 'react'
import {render} from 'react-dom'
import reportWebVitals from './reportWebVitals'
import {Provider} from 'react-redux'
import configureStore from "./store"
import {getConfigAsync} from "./features/config/configSlice"
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom"
import AdminLayout from "./layouts/Admin.js"

import "bootstrap/dist/css/bootstrap.min.css"
import "./assets/css/animate.min.css"
import "./assets/scss/light-bootstrap-dashboard-react.scss?v=2.0.0"
import "./assets/css/demo.css"
import "@fortawesome/fontawesome-free/css/all.min.css"

const store = configureStore()
store.dispatch(getConfigAsync())

const renderApp = () =>
  render(
    <Provider store={store}>
      <BrowserRouter>
        <Switch>
          <Route path="/admin" render={(props) => <AdminLayout {...props} />}/>
          <Redirect from="/" to="/admin/config"/>
        </Switch>
      </BrowserRouter>
    </Provider>,
    document.getElementById('root')
  )

renderApp()
// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
