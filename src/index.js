import React, { useState, useEffect } from "react";
import { render } from "react-dom";
import "./index.css";

const fetchLinks = () => fetch("/link").then(r => r.json());

const addLink = link =>
  fetch("/link", { method: "POST", body: link }).then(r => r.json());

const Shortner = () => {
  const [links, setLinks] = useState([]);
  useEffect(() => {
    fetchLinks().then(setLinks);
  }, []);
  return (
    <div className="container">
      <h1>Shortner</h1>
      <InputForm setLinks={setLinks} />
      <Table links={links} />
    </div>
  );
};

const InputForm = ({ setLinks }) => {
  const [link, setLink] = useState("");
  return (
    <div className="form">
      <input
        className="form-control"
        value={link}
        onChange={e => setLink(e.target.value)}
        type="url"
      />

      <button
        className="btn btn-default"
        onClick={() => {
          if (!link) return;
          setLink("");
          addLink(link).then(setLinks);
        }}
      >
        Сократить
      </button>
    </div>
  );
};
const absolute = path => `${window.location.href}${path}`;
const Table = ({ links }) => (
  <table className="table table-striped table-bordered table-condensed">
    <thead>
      <tr>
        <th>Короткая</th>
        <th>Исходная</th>
      </tr>
    </thead>
    <tbody>
      {links.map(({ Long, Short }) => (
        <tr key={Short}>
          <td>
            <a href={absolute(Short)}>{absolute(Short)}</a>
          </td>
          <td>
            <a href={Long}>{Long.replace(/^(.{0,100})(.*)/g, "$1...")}</a>
          </td>
        </tr>
      ))}
    </tbody>
  </table>
);

render(<Shortner />, document.getElementById("root"));
