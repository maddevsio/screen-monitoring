import React from 'react';
import Pages from './components/Pages.jsx'

const INTERVAL = 1000

class App extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            pages: props.pages || []
        };
        this.receiveResponse = (data) => {
            this.setState({
                pages: data.widgets
            });
        };

        this._renderPages = () => {
            if(this.state.pages.length > 0) {
                return <Pages pages={this.state.pages}/>
            }else {
                return (
                    <strong>No widgets</strong>
                )
            }
        }

        this._grabInfo = () => {
             var self = this;
             fetch('/dashboard/v1/pages')
                         .then((response) => response.json())
                         .then(self.receiveResponse)
                         .catch((ex) => {
                             console.error(ex);
                             setTimeout(self._grabInfo, INTERVAL);
                         });
        }
    }

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
