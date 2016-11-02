import React from 'react';
import { Router, Route, hashHistory } from 'react-router';
import {render} from 'react-dom';
import App from './App.jsx';
import UnregisteredWidgetList from './components/UnregisteredWidgetList.jsx';



render((
  <Router history={hashHistory}>
    <Route path="/" component={App}/>
    <Route path="/unregistered" component={UnregisteredWidgetList}/>
  </Router>
),
document.getElementById('app'));
