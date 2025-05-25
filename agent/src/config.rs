use goldberg::goldberg_string;

#[inline]
pub fn get_address() -> String {
    // Put your onion domain here!!!
    goldberg_string!("").to_string()
}

#[inline]
pub fn get_id_path() -> String {
    goldberg_string!("agent.id").to_string()
}