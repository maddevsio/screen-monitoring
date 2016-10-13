import React from 'react';
import Preview from './Preview.jsx';

class Pages extends React.Component {
    static defaultProps = {
        pages: [],
        pageIndex: 0
    }

    constructor(props) {
        super(props);
        this.state = {
            pages: props.pages || [],
            pageIndex: props.pageIndex
        };
    }

    renderPage = () => {
        var page = this.state.pages[this.state.pageIndex];
        return page.content.map((w, idx) => {
            return <Preview {...w} key={w.id} />
        });
    }

    componentWillReceiveProps(nextProps) {
      this.setState({
        pages: nextProps.pages || [],
        pageIndex: nextProps.pageIndex
      });
    }

    render() {
        return (
            <div>
                {this.renderPage()}
            </div>
        );
    }
}

export default Pages
