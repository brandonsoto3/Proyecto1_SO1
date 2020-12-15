//INCLUIMOS LAS BIBLIOTECAS NECESARIAS
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>

//DEFINIENDO EL TAMAÑO DEL BUFFER
#define BUFSIZE  	150

//DEFINIENDO INFORMACIÓN ACERCA DEL MODULO
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo que escribe informacion acerca de la memoria RAM");
MODULE_AUTHOR("BRANDON JAVIER SOTO CASTAÑEDA 201503893");



static int iniciar(void)
{
    proc_create("memo_201503893", 0, NULL, &operaciones);
    printk(KERN_INFO "201503893\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("memo_201503893", NULL);
    printk(KERN_INFO "Sistemas Operativos 1 \n");
}
 
module_init(iniciar);
module_exit(salir); 


