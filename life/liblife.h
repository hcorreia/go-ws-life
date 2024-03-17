void *init_state_random(int width, int height, int n_workers);

void *next_state(void *game_ptr);

void *next_state_img(void *game_ptr);
// char *next_state_img(void *game_ptr);

void free_char_p(void *s);

void free_void_p(void *ptr);
