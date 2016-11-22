import React from 'react';
import { Router, Route, hashHistory } from 'react-router';
import {render} from 'react-dom';
import App from './App.jsx';
import UnregisteredWidgetList from './components/UnregisteredWidgetList.jsx';
import AddNewPageContainer from './components/AddNewPageContainer.jsx';
import PagesListContainer from './components/PagesListContainer.jsx';

render((
  <Router history={hashHistory}>
    <Route path="/" component={App}/>
    <Route path="/unregistered" component={UnregisteredWidgetList}/>
    <Route path="/pages/list" component={PagesListContainer}/>
    <Route path="/pages/list/:selectedWidgetId/" component={PagesListContainer}/>
    <Route path="/pages/new" component={AddNewPageContainer}/>
    <Route path="/pages/new/:selectedWidgetId/" component={AddNewPageContainer}/>
  </Router>
),
document.getElementById('app'));
