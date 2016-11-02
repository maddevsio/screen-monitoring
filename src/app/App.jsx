import React, { PropTypes, PureComponent } from 'react';
import Pages from './components/Pages.jsx'

const INTERVAL = 1000

class App extends PureComponent {

    constructor(props) {
        super(props);
        this.state = {
            pages: props.pages || []
        };
    }

    _grabInfo = () => {
         var self = this;
         fetch('/dashboard/v1/pages')
                     .then((response) => response.json())
                     .then(self._receiveResponse)
                     /*
                     .then(function (){
                        setTimeout(self._grabInfo, INTERVAL);
                     })
                     */
                     .catch((ex) => {
                         console.error(ex);
                         setTimeout(self._grabInfo, INTERVAL);
                     });
    }

    _renderPages = () => {
        if(this.state.pages.length > 0) {
            return <Pages pages={this.state.pages}/>
        }else {
            return (
                <strong>No widgets</strong>
            )
        }
    }

    _receiveResponse = (data) => {
        this.setState({
            pages: data
        });
    };

    componentDidMount(){
        this._grabInfo();
    }

    render() {
        return (
            <div style={{width:'100%', height:'100%'}}>
                {this._renderPages()}
            </div>
        )

    }
}

export default App;
