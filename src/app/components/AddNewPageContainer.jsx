import React, {PureComponent} from 'react';
import { hashHistory } from 'react-router';
import AddPageForm from './AddPageForm.jsx';

class AddNewPageContainer extends PureComponent {

  _onCreatePage = (p) => {
    var self = this;
    const CREATE_PAGE_URL = '/dashboard/v1/pages/new';
    fetch(CREATE_PAGE_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(p)
    })
    .then(r => r.json())
    .then(data => {
       if (data.Success) {
         var selectedId = this.props.location.state.selectedWidgetId;
         hashHistory.push({
           pathname: '/pages/list',
           state: { selectedWidgetId: selectedId }
         });
       } else {
         alert(data.error);
       }
    })
    .catch((ex) => {
        console.error(ex);
    });
  }

  _onCancelPage = () => {
    var selectedId = this.props.location.state.selectedWidgetId;
    hashHistory.push({
      pathname: '/pages/list',
      state: { selectedWidgetId: selectedId }
    });
  }

  render() {
    return (
      <AddPageForm onCreate={this._onCreatePage} onCancel={this._onCancelPage}/>
    )
  }
}

export default AddNewPageContainer;
