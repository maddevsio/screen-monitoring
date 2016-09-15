import React from 'react';
import Preview from './Preview.jsx';

class Pages extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            pages: props.pages || []
        };

        this.renderPages = () => {
               return this.state.pages.map((page, index) => {
                   return <Preview {...page} key={index} />
               })
        }
    }

    componentWillReceiveProps(nextProps) {
      this.setState({
        pages: nextProps.pages || []
      });
    }

    render() {
        return (
            <div>
                {this.renderPages()}
            </div>
        );
    }
}

export default Pages