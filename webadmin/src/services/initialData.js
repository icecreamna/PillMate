import { request } from "../lib/http";

/** ===== Forms ===== */
// GET /forms[?with_relations=true]
export const listForms = ({ with_relations = false } = {}) =>
  request(`/forms${with_relations ? "?with_relations=true" : ""}`); // -> array

// GET /form/:id[?with_relations=true]
export const getForm = (id, { with_relations = false } = {}) =>
  request(`/form/${id}${with_relations ? "?with_relations=true" : ""}`); // -> object

// GET /forms/:id/units  -> { form_id, units: [{id, unit_name}] }
export const listUnitsByForm = (formId) =>
  request(`/forms/${formId}/units`);

/** ===== Units ===== */
// GET /units[?with_relations=true]
export const listUnits = ({ with_relations = false } = {}) =>
  request(`/units${with_relations ? "?with_relations=true" : ""}`);

// GET /unit/:id[?with_relations=true]
export const getUnit = (id, { with_relations = false } = {}) =>
  request(`/unit/${id}${with_relations ? "?with_relations=true" : ""}`);

// GET /units/:id/forms -> { unit_id, forms: [{id, form_name}] }
export const listFormsByUnit = (unitId) =>
  request(`/units/${unitId}/forms`);

/** ===== Instructions ===== */
// GET /instructions
export const listInstructions = () => request(`/instructions`);

// GET /instruction/:id
export const getInstruction = (id) => request(`/instruction/${id}`);
