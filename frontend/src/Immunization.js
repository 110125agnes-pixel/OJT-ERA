import React, { useState } from 'react';
import './Immunization.css';

function Checkbox({ label, checked, onChange }) {
  return (
    <label className="immu-checkbox">
      <input type="checkbox" checked={checked} onChange={onChange} />
      <span>{label}</span>
    </label>
  );
}

export default function Immunization() {
  // Local-only state; no backend calls as requested
  const [children, setChildren] = useState({
    none: false,
    bcg: true,
    opv1: true,
    opv2: false,
    opv3: false,
    dpt1: false,
    dpt2: false,
    dpt3: true,
    measles: false,
    hepB1: false,
    hepB2: false,
    hepB3: false,
    hepA: false,
    varicella: true,
  });

  const [young, setYoung] = useState({ none: false, hpv: false, mmr: true });
  const [pregnant, setPregnant] = useState({ none: true, tetanus: false });
  const [elderly, setElderly] = useState({ none: true, pneumo: false, flu: false });
  const [other, setOther] = useState('lorem ipsum 12345');

  const toggle = (groupSetter, group, key) => {
    groupSetter({ ...group, [key]: !group[key] });
  };

  return (
    <div className="immu-container">
      <div className="immu-header">
        <h2>Immunization</h2>
        <div className="immu-tabs">
          <span className="immu-tab active">Profiling</span>
          <span className="immu-tab">SOAP</span>
          <span className="immu-tab">Medicine</span>
        </div>
      </div>

      <div className="immu-sections">
        <section className="immu-section">
          <h3>1. Children</h3>
          <div className="immu-grid">
            <Checkbox label="None" checked={children.none} onChange={() => toggle(setChildren, children, 'none')} />
            <Checkbox label="BCG" checked={children.bcg} onChange={() => toggle(setChildren, children, 'bcg')} />
            <Checkbox label="OPV1" checked={children.opv1} onChange={() => toggle(setChildren, children, 'opv1')} />
            <Checkbox label="OPV2" checked={children.opv2} onChange={() => toggle(setChildren, children, 'opv2')} />
            <Checkbox label="OPV3" checked={children.opv3} onChange={() => toggle(setChildren, children, 'opv3')} />
            <Checkbox label="DPT1" checked={children.dpt1} onChange={() => toggle(setChildren, children, 'dpt1')} />
            <Checkbox label="DPT2" checked={children.dpt2} onChange={() => toggle(setChildren, children, 'dpt2')} />
            <Checkbox label="DPT3" checked={children.dpt3} onChange={() => toggle(setChildren, children, 'dpt3')} />
            <Checkbox label="Measles" checked={children.measles} onChange={() => toggle(setChildren, children, 'measles')} />
            <Checkbox label="Hepatitis B1" checked={children.hepB1} onChange={() => toggle(setChildren, children, 'hepB1')} />
            <Checkbox label="Hepatitis B2" checked={children.hepB2} onChange={() => toggle(setChildren, children, 'hepB2')} />
            <Checkbox label="Hepatitis B3" checked={children.hepB3} onChange={() => toggle(setChildren, children, 'hepB3')} />
            <Checkbox label="Hepatitis A" checked={children.hepA} onChange={() => toggle(setChildren, children, 'hepA')} />
            <Checkbox label="Varicella (Chicken Pox)" checked={children.varicella} onChange={() => toggle(setChildren, children, 'varicella')} />
          </div>
        </section>

        <section className="immu-section">
          <h3>2. Young</h3>
          <div className="immu-grid">
            <Checkbox label="None" checked={young.none} onChange={() => toggle(setYoung, young, 'none')} />
            <Checkbox label="HPV" checked={young.hpv} onChange={() => toggle(setYoung, young, 'hpv')} />
            <Checkbox label="MMR" checked={young.mmr} onChange={() => toggle(setYoung, young, 'mmr')} />
          </div>
        </section>

        <section className="immu-section">
          <h3>3. Pregnant</h3>
          <div className="immu-grid">
            <Checkbox label="None" checked={pregnant.none} onChange={() => toggle(setPregnant, pregnant, 'none')} />
            <Checkbox label="Tetanus Toxoid" checked={pregnant.tetanus} onChange={() => toggle(setPregnant, pregnant, 'tetanus')} />
          </div>
        </section>

        <section className="immu-section">
          <h3>4. Elderly</h3>
          <div className="immu-grid">
            <Checkbox label="None" checked={elderly.none} onChange={() => toggle(setElderly, elderly, 'none')} />
            <Checkbox label="Pneumococcal Vaccine" checked={elderly.pneumo} onChange={() => toggle(setElderly, elderly, 'pneumo')} />
            <Checkbox label="Flu Vaccine" checked={elderly.flu} onChange={() => toggle(setElderly, elderly, 'flu')} />
          </div>
        </section>

        <section className="immu-section">
          <h3>5. Other Immunization</h3>
          <textarea
            className="immu-notes"
            value={other}
            onChange={(e) => setOther(e.target.value)}
            placeholder="Add notes here"
          />
        </section>

        <div className="immu-footer">
          <div className="immu-datetime">
            <input type="date" />
            <input type="time" />
          </div>
          <button className="immu-close">Close</button>
        </div>
      </div>
    </div>
  );
}
