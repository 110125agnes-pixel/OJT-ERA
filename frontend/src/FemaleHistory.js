import React, { useState } from 'react';
import './FemaleHistory.css';

function NumberSelect({ value, onChange, min = 0, max = 20 }) {
  const options = [];
  for (let i = min; i <= max; i += 1) options.push(i);
  return (
    <select value={value} onChange={(e) => onChange(Number(e.target.value))}>
      {options.map((n) => (
        <option key={n} value={n}>
          {n}
        </option>
      ))}
    </select>
  );
}

function YesNo({ value, onChange, name }) {
  return (
    <div className="fh-radio-group">
      <label>
        <input
          type="radio"
          name={name}
          checked={value === true}
          onChange={() => onChange(true)}
        />
        Yes
      </label>
      <label>
        <input
          type="radio"
          name={name}
          checked={value === false}
          onChange={() => onChange(false)}
        />
        No
      </label>
    </div>
  );
}

export default function FemaleHistory() {
  // Menstrual state
  const [menstrual, setMenstrual] = useState({
    firstAge: 12,
    lastPeriod: '2025-11-14',
    durationDays: 7,
    cycleDays: 30,
    padsPerDay: 4,
    intercourseAge: 18,
    birthControl: 'ORA',
    menopause: false,
    menopauseAge: 40,
    applicable: false,
  });

  // Pregnancy state
  const [preg, setPreg] = useState({
    gravity: 4,
    parity: 2,
    deliveryType: 'normal',
    fullTerm: 0,
    premature: 0,
    abortion: 0,
    livingChildren: 3,
    eclampsia: true,
    familyPlanning: true,
    applicable: false,
  });

  const updateMenstrual = (key, value) => setMenstrual({ ...menstrual, [key]: value });
  const updatePreg = (key, value) => setPreg({ ...preg, [key]: value });

  return (
    <div className="fh-container">
      <div className="fh-header">
        <div className="fh-tabs">
          <span className="fh-tab active">Profiling</span>
          <span className="fh-tab">SOAP</span>
          <span className="fh-tab">Medicine</span>
        </div>
        <div className="fh-subtabs">
          <span className="fh-subtab">Medical</span>
          <span className="fh-subtab">Surgery</span>
          <span className="fh-subtab">Family</span>
          <span className="fh-subtab">Immunization</span>
          <span className="fh-subtab">Social History</span>
          <span className="fh-subtab active">Female</span>
          <span className="fh-subtab">*Pertinent Physical Examinations</span>
          <span className="fh-subtab">*Physical Examination</span>
          <span className="fh-subtab">NCDQANS</span>
        </div>
      </div>

      <div className="fh-panels">
        <section className="fh-panel">
          <div className="fh-panel-title">Menstrual History</div>
          <div className="fh-grid">
            <label className="fh-row">
              <span>1. Age of First Menstruation (Menarche)</span>
              <NumberSelect value={menstrual.firstAge} onChange={(v) => updateMenstrual('firstAge', v)} min={8} max={20} />
            </label>
            <label className="fh-row">
              <span>2. Date of Last Menstrual Period</span>
              <input
                type="date"
                value={menstrual.lastPeriod}
                onChange={(e) => updateMenstrual('lastPeriod', e.target.value)}
              />
            </label>
            <label className="fh-row">
              <span>3. Duration of Menstrual Period in Number of Days</span>
              <NumberSelect value={menstrual.durationDays} onChange={(v) => updateMenstrual('durationDays', v)} min={1} max={15} />
            </label>
            <label className="fh-row">
              <span>4. Interval/Cycle of Menstruation in Number of Days</span>
              <NumberSelect value={menstrual.cycleDays} onChange={(v) => updateMenstrual('cycleDays', v)} min={10} max={60} />
            </label>
            <label className="fh-row">
              <span>5. Number of Pads/Napkins Used per Day during Menstruation</span>
              <NumberSelect value={menstrual.padsPerDay} onChange={(v) => updateMenstrual('padsPerDay', v)} min={0} max={15} />
            </label>
            <label className="fh-row">
              <span>6. Onset of Sexual Intercourse (Age of First Sexual Intercourse)</span>
              <NumberSelect value={menstrual.intercourseAge} onChange={(v) => updateMenstrual('intercourseAge', v)} min={10} max={40} />
            </label>
            <label className="fh-row">
              <span>7. Birth Control Method Used</span>
              <input
                type="text"
                value={menstrual.birthControl}
                onChange={(e) => updateMenstrual('birthControl', e.target.value)}
                placeholder="Enter method"
              />
            </label>
            <label className="fh-row">
              <span>8. Is Menopause?</span>
              <YesNo value={menstrual.menopause} onChange={(v) => updateMenstrual('menopause', v)} name="fh-menopause" />
            </label>
            <label className="fh-row">
              <span>9. If Menopause, Age of Menopause</span>
              <NumberSelect value={menstrual.menopauseAge} onChange={(v) => updateMenstrual('menopauseAge', v)} min={20} max={70} />
            </label>
            <label className="fh-row">
              <span>10. Is menstrual history applicable?</span>
              <YesNo value={menstrual.applicable} onChange={(v) => updateMenstrual('applicable', v)} name="fh-menstrual-app" />
            </label>
          </div>
        </section>

        <section className="fh-panel">
          <div className="fh-panel-title">Pregnancy History</div>
          <div className="fh-grid">
            <label className="fh-row">
              <span>1. Number of Pregnancy to Date – Gravity Chief</span>
              <NumberSelect value={preg.gravity} onChange={(v) => updatePreg('gravity', v)} min={0} max={20} />
            </label>
            <label className="fh-row">
              <span>2. Number of Delivery to Date – Parity</span>
              <NumberSelect value={preg.parity} onChange={(v) => updatePreg('parity', v)} min={0} max={20} />
            </label>
            <label className="fh-row">
              <span>3. Type of Delivery</span>
              <div className="fh-radio-group">
                <label>
                  <input
                    type="radio"
                    name="fh-delivery"
                    checked={preg.deliveryType === 'normal'}
                    onChange={() => updatePreg('deliveryType', 'normal')}
                  />
                  Normal
                </label>
                <label>
                  <input
                    type="radio"
                    name="fh-delivery"
                    checked={preg.deliveryType === 'operative'}
                    onChange={() => updatePreg('deliveryType', 'operative')}
                  />
                  Operative
                </label>
                <label>
                  <input
                    type="radio"
                    name="fh-delivery"
                    checked={preg.deliveryType === 'both'}
                    onChange={() => updatePreg('deliveryType', 'both')}
                  />
                  Both Normal and Operative
                </label>
                <label>
                  <input
                    type="radio"
                    name="fh-delivery"
                    checked={preg.deliveryType === 'na'}
                    onChange={() => updatePreg('deliveryType', 'na')}
                  />
                  Not Applicable
                </label>
              </div>
            </label>
            <label className="fh-row">
              <span>4. Number of Full Term Pregnancy</span>
              <NumberSelect value={preg.fullTerm} onChange={(v) => updatePreg('fullTerm', v)} min={0} max={20} />
            </label>
            <label className="fh-row">
              <span>5. Number of Premature Pregnancy</span>
              <NumberSelect value={preg.premature} onChange={(v) => updatePreg('premature', v)} min={0} max={20} />
            </label>
            <label className="fh-row">
              <span>6. Number of Abortion</span>
              <NumberSelect value={preg.abortion} onChange={(v) => updatePreg('abortion', v)} min={0} max={20} />
            </label>
            <label className="fh-row">
              <span>7. Number of Living Children</span>
              <NumberSelect value={preg.livingChildren} onChange={(v) => updatePreg('livingChildren', v)} min={0} max={20} />
            </label>
            <label className="fh-row">
              <span>8. If Pregnancy – Induced Hypertension (Pre – Eclampsia)</span>
              <YesNo value={preg.eclampsia} onChange={(v) => updatePreg('eclampsia', v)} name="fh-eclampsia" />
            </label>
            <label className="fh-row">
              <span>9. If with access to Family Planning Counselling</span>
              <YesNo value={preg.familyPlanning} onChange={(v) => updatePreg('familyPlanning', v)} name="fh-familyplan" />
            </label>
            <label className="fh-row">
              <span>10. Is pregnancy history applicable?</span>
              <YesNo value={preg.applicable} onChange={(v) => updatePreg('applicable', v)} name="fh-preg-app" />
            </label>
          </div>
        </section>

        <div className="fh-footer">
          <div className="fh-datetime">
            <input type="date" />
            <input type="time" />
          </div>
          <button className="fh-close">Close</button>
        </div>
      </div>
    </div>
  );
}
