#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/sched.h>        
#include <linux/sched/signal.h> 

//DEFINIENDO EL TAMAÑO DEL BUFFER
#define BUFSIZE  	150

//DEFINIENDO INFORMACIÓN ACERCA DEL MODULO
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo que escribe informacion acerca del CPU");
MODULE_AUTHOR("BRANDON JAVIER SOTO CASTAÑEDA 201503893");

struct sysinfo inf;

static int iniciar(void)
{
    proc_create("cpu_201503893", 0, NULL, &operaciones);
    printk(KERN_INFO "BRANDON JAVIER SOTO CASTANEDA\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("cpu_201503893", NULL);
    printk(KERN_INFO "Diciembre 2020\n");
}
 
module_init(iniciar);
module_exit(salir); 