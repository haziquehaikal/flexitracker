import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { Link } from "react-router-dom";

import Login from "./Login";
import Home from "./Home";
import Product from "./Product";
// import HistoryData from '../history.xml';
import HistoryPage from "./Login";
// component={() => <Site tileUrl2={process.env.TILE_URL}

const App = () => (
  <Router>
    <Switch>
      <Route exact path="/" component={Home} />
      <Route path="/product" component={Product} />
      <Route path="/analytic" component={HistoryPage} />
      <Route path="/staff" component={HistoryPage} />
    </Switch>
  </Router>
);

export default App;
