
#ifndef RUBY_GC_H
#define RUBY_GC_H 1

#if defined(__x86_64__) && defined(__GNUC__)
#define SET_MACHINE_STACK_END(p) __asm__("movq\t%%rsp, %0" : "=r" (*p))
#elif defined(__i386) && defined(__GNUC__)
#define SET_MACHINE_STACK_END(p) __asm__("movl\t%%esp, %0" : "=r" (*p))
#else
NOINLINE(void rb_gc_set_stack_end(VALUE **stack_end_p));
#define SET_MACHINE_STACK_END(p) rb_gc_set_stack_end(p)
#define USE_CONSERVATIVE_STACK_END
#endif

/* for GC debug */

#ifndef RUBY_MARK_FREE_DEBUG
#define RUBY_MARK_FREE_DEBUG 0
#endif

#if RUBY_MARK_FREE_DEBUG
extern int ruby_gc_debug_indent;

static void
rb_gc_debug_indent(void)
{
    printf("%*s", ruby_gc_debug_indent, "");
}

static void
rb_gc_debug_body(const char *mode, const char *msg, int st, void *ptr)
{
    if (st == 0) {
	ruby_gc_debug_indent--;
    }
    rb_gc_debug_indent();
    printf("%s: %s %s (%p)\n", mode, st ? "->" : "<-", msg, ptr);

    if (st) {
	ruby_gc_debug_indent++;
    }

    fflush(stdout);
}

#define RUBY_MARK_ENTER(msg) rb_gc_debug_body("mark", msg, 1, ptr)
#define RUBY_MARK_LEAVE(msg) rb_gc_debug_body("mark", msg, 0, ptr)
#define RUBY_FREE_ENTER(msg) rb_gc_debug_body("free", msg, 1, ptr)
#define RUBY_FREE_LEAVE(msg) rb_gc_debug_body("free", msg, 0, ptr)
#define RUBY_GC_INFO         rb_gc_debug_indent(); printf

#else
#define RUBY_MARK_ENTER(msg)
#define RUBY_MARK_LEAVE(msg)
#define RUBY_FREE_ENTER(msg)
#define RUBY_FREE_LEAVE(msg)
#define RUBY_GC_INFO if(0)printf
#endif

#define RUBY_MARK_UNLESS_NULL(ptr) if(RTEST(ptr)){rb_gc_mark(ptr);}
#define RUBY_FREE_UNLESS_NULL(ptr) if(ptr){ruby_xfree(ptr);}

#if STACK_GROW_DIRECTION > 0
# define STACK_UPPER(x, a, b) a
#elif STACK_GROW_DIRECTION < 0
# define STACK_UPPER(x, a, b) b
#else
RUBY_EXTERN int ruby_stack_grow_direction;
int ruby_get_stack_grow_direction(volatile VALUE *addr);
# define stack_growup_p(x) (			\
	(ruby_stack_grow_direction ?		\
	 ruby_stack_grow_direction :		\
	 ruby_get_stack_grow_direction(x)) > 0)
# define STACK_UPPER(x, a, b) (stack_growup_p(x) ? a : b)
#endif

#endif /* RUBY_GC_H */
